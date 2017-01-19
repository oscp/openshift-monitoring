package client

import (
	"github.com/valyala/gorpc"
	"log"
)

var c *gorpc.Client

func RegisterOnHub(hubAddr string) {
	log.Println("trying to contact hub on: ", hubAddr)

	// Register on hub
	c = &gorpc.Client{
		Addr: hubAddr,
	}
	c.Start()

	resp, err := c.Call("foobar")
	if err != nil {
		log.Fatalf("Error when sending request to hub: %s", err)
	}
	if resp.(string) != "foobar" {
		log.Fatalf("Unexpected response from the hub: %+v", resp)
	}
}

func DeamonServer(port string) {
	log.Println("creating deamon server on port: ", port)
	s := &gorpc.Server{
		Addr: ":" + port,
		Handler: func(clientAddr string, request interface{}) interface{} {
			log.Printf("new job from hub: ", request)
			return request
		},
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("Cannot start deamon server: %s", err)
	}
}
