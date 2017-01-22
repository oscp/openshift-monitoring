package server

import (
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"github.com/cenkalti/rpc2"
	"net"
)

type Hub struct {
	hubAddr   string
	deamons   map[string]models.DeamonClient
	jobs	  []models.Job
	lastJobId int64
	toDeamons chan models.Job
	toUi      chan models.BaseModel
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		hubAddr: hubAddr,
		deamons: make(map[string]models.DeamonClient),
		jobs: []models.Job{},
		lastJobId: 0,
		toDeamons: make(chan models.Job),
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
	go handleToDeamons(h)

	srv := rpc2.NewServer()
	srv.Handle("register", func(c *rpc2.Client, d *models.Deamon, reply *string) error {

		// Save client for talking to him
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

func handleToDeamons(h *Hub) {
	log.Println("ready to send jobs to deamons")
	for {
		var job models.Job = <- h.toDeamons

		for _,d := range h.deamons {
			if err := d.Client.Call("startJob", job, nil); err != nil {
				log.Println("error starting job on deamon", err)
			}
		}
	}

}
