package client

import (
	"log"
	"net"
	"github.com/cenkalti/rpc2"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"os"
)

func RegisterOnHub(h string, dt string) *rpc2.Client {
	// Register on hub
	conn, _ := net.Dial("tcp", h)
	c := rpc2.NewClient(conn)
	c.Handle("startChecks", func(client *rpc2.Client, checks *models.Checks, reply *string) error {
		startChecks(checks)
		*reply = "ok"
		return nil
	})
	c.Handle("stopChecks", func(client *rpc2.Client, stop *bool, reply *string) error {
		stopChecks()
		*reply = "ok"
		return nil
	})
	go c.Run()

	// Register on hub
	log.Println("registring on hub:", h)
	var rep string
	host, _ := os.Hostname()
	err := c.Call("register", models.Deamon{Hostname: host, DeamonType: dt, ChecksCount: 0}, &rep)
	if err != nil {
		log.Fatal("error registring on hub: ", err)
	}
	if rep != "ok" {
		log.Fatalf("expected the hub to answer with ok. he did with: %+v", rep)
	}

	return c
}

func UnregisterOnHub(c *rpc2.Client) {
	log.Println("unregistring from hub")

	var rep string
	host, _ := os.Hostname()
	err := c.Call("unregister", host, &rep)
	if err != nil {
		log.Fatalf("error when unregistring from hub: %s", err)
	}
	c.Close()
}

