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
	toDeamons chan []byte
	toUi      chan models.BaseModel
	uiLeave   chan struct{}
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		deamons: make(map[string]models.DeamonClient),
		toDeamons: make(chan []byte),
		toUi: make(chan models.BaseModel, 1000),
		hubAddr: hubAddr,
		uiLeave: make(chan struct{}),
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
