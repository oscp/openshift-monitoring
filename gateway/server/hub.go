package server

import (
	"github.com/valyala/gorpc"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
)

type Hub struct {
	hubAddr string
	deamons   map[string]models.Deamon
	toDeamons chan []byte
	toUi      chan models.BaseModel
}

func NewHub(hubAddr string) *Hub {
	return &Hub{
		deamons: make(map[string]models.Deamon),
		toDeamons: make(chan []byte),
		toUi: make(chan models.BaseModel, 1000),
		hubAddr: hubAddr,
	}
}

func (h *Hub) Serve() {
	// register models
	gorpc.RegisterType(&models.Deamon{})

	s := &gorpc.Server{
		Addr: h.hubAddr,
		Handler: func(clientAddr string, r interface{}) interface{} {
			switch v := r.(type) {
			case *models.Deamon:
				deamonJoin(h, clientAddr, v)
				break;
			case string:
				deamonLeave(h, clientAddr)
				break;
			default:
				log.Println("unknown type on rpc ", r, v)
			}

			return "ok"
		},
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}
}
