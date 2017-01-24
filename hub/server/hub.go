package server

import (
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"github.com/cenkalti/rpc2"
	"net"
)

type Hub struct {
	hubAddr       string
	deamons       map[string]models.DeamonClient
	currentChecks models.Checks
	startChecks   chan models.Checks
	stopChecks    chan bool
	toUi          chan models.BaseModel
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		hubAddr: hubAddr,
		deamons: make(map[string]models.DeamonClient),
		startChecks: make(chan models.Checks),
		stopChecks: make(chan bool),
		toUi: make(chan models.BaseModel, 1000),
	}
}

func (h *Hub) Deamons() []models.Deamon {
	r := []models.Deamon{}
	for _, d := range h.deamons {
		r = append(r, d.Deamon)
	}
	return r
}

func (h *Hub) Serve() {
	go handleChecksStart(h)
	go handleChecksStop(h)

	srv := rpc2.NewServer()
	srv.Handle("register", func(c *rpc2.Client, d *models.Deamon, reply *string) error {
		// Save client for talking to him later
		deamonJoin(h, d, c)

		*reply = "ok"
		return nil
	})
	srv.Handle("unregister", func(cl *rpc2.Client, host *string, reply *string) error {
		deamonLeave(h, *host)

		*reply = "ok"
		return nil
	})

	lis, err := net.Listen("tcp", h.hubAddr)
	srv.Accept(lis)
	if err != nil {
		log.Fatalf("Cannot start rpc2 server: %s", err)
	}
}

func handleChecksStart(h *Hub) {
	for {
		var checks models.Checks = <-h.startChecks
		log.Println("Sending to deamons", checks)

		for _, d := range h.deamons {
			if err := d.Client.Call("startChecks", checks, nil); err != nil {
				log.Println("error starting checks on deamon", err)
			}
		}
	}
}

func handleChecksStop(h *Hub) {
	for {
		var stop bool = <-h.stopChecks

		if (stop) {
			log.Println("Sending stop command to deamons")
			for _, d := range h.deamons {
				if err := d.Client.Call("stopChecks", stop, nil); err != nil {
					log.Println("error stopping checks on deamon", err)
				}
			}
		}
	}
}
