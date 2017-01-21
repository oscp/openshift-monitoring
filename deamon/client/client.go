package client

import (
	"github.com/valyala/gorpc"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"strconv"
)

var c *gorpc.Client

func RegisterOnHub(h string, dt string, p int) *gorpc.Client {
	log.Println("trying to contact hub on: ", h)

	// Register on hub
	gorpc.RegisterType(&models.Deamon{})
	c = &gorpc.Client{
		Addr: h,
	}
	c.Start()

	resp, err := c.Call(models.Deamon{Addr: h, DeamonType: dt, Port: p})
	if err != nil {
		log.Fatalf("error when sending request to hub: %s", err)
	}
	if resp.(string) != "ok" {
		log.Fatalf("expected the hub to answer with ok. he did not: %+v", resp)
	}

	return c
}

func UnregisterOnHub(c *gorpc.Client) {
	log.Println("unregistring from hub")

	_, err := c.Call("unregister")
	if err != nil {
		log.Fatalf("error when unregistring from hub: %s", err)
	}
	c.Stop()
}

func DeamonServer(p int) {
	log.Println("creating deamon server on port: ", p)
	s := &gorpc.Server{
		Addr: ":" + strconv.Itoa(p),
		Handler: func(clientAddr string, request interface{}) interface{} {
			log.Printf("new job from hub: ", request)
			return request
		},
	}
	if err := s.Serve(); err != nil {
		log.Fatalf("cannot start deamon server: %s", err)
	}
}
