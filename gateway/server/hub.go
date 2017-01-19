package server

import (
	"github.com/valyala/gorpc"
	"log"
)

type Deamon struct {
	addr string
}

type Hub struct {
	deamons []Deamon
}

var hub = &Hub{
	deamons: []Deamon{},
}

func DeamonHub(hubAddr string) {
	s := &gorpc.Server{
		Addr: hubAddr,
		Handler: func(clientAddr string, request interface{}) interface{} {
			log.Printf("new deamon joined, %+v, %s\n", request, clientAddr)
			hub.deamons = append(hub.deamons, Deamon{ addr: clientAddr })
			return request
		},
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start rpc server: %s", err)
	}
}
