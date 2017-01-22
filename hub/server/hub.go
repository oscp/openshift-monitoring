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
	jobs      map[int64]*models.Job
	lastJobId int64
	jobStart  chan models.Job
	jobStop   chan int64
	toUi      chan models.BaseModel
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		hubAddr: hubAddr,
		deamons: make(map[string]models.DeamonClient),
		jobs: make(map[int64]*models.Job),
		lastJobId: 0,
		jobStart: make(chan models.Job),
		jobStop: make(chan int64),
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

func (h *Hub) Jobs() []models.Job {
	r := []models.Job{}
	for _,j := range h.jobs {
		r = append(r, *j)
	}
	return r
}

func (h *Hub) Serve() {
	go handleJobStart(h)
	go handleJobStop(h)

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

func handleJobStart(h *Hub) {
	for {
		var job models.Job = <- h.jobStart

		for _,d := range h.deamons {
			if err := d.Client.Call("startJob", job, nil); err != nil {
				log.Println("error starting job on deamon", err)
			}
		}
	}
}

func handleJobStop(h *Hub) {
	for {
		var jobId int64 = <- h.jobStop

		for _,d := range h.deamons {
			if err := d.Client.Call("stopJob", jobId, nil); err != nil {
				log.Println("error stopping job on deamon", err)
			}
		}
	}

}
