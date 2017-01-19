package server

import (
	"github.com/valyala/gorpc"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
)

type Deamon struct {
	addr string
}

type Hub struct {
	hubAddr string

	deamons   []Deamon

	// send things to deamons
	toDeamons chan []byte

	// send things to ui
	toUi      chan models.BaseModel
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		deamons: []Deamon{},
		toDeamons: make(chan []byte),
		toUi: make(chan models.BaseModel, 1000),
		hubAddr: hubAddr,
	}
}

func (h *Hub) Serve() {
	s := &gorpc.Server{
		Addr: h.hubAddr,
		Handler: func(clientAddr string, request interface{}) interface{} {
			log.Printf("new deamon joined, %+v, %s\n", request, clientAddr)
			h.deamons = append(h.deamons, Deamon{addr: clientAddr })

			// tell the ui about it
			h.toUi <- models.BaseModel{ Type: models.TYPE_NEW_DEAMON, Message: clientAddr }
			return "ok"
		},
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}
}
