package client

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"log"
	"os"
	"github.com/cenkalti/rpc2"
)

func registerOnHub(h string, dc *models.DeamonClient) {
	log.Println("registring on hub:", h)
	var rep string
	err := dc.Client.Call("register", dc.Deamon, &rep)
	if err != nil {
		log.Fatal("error registring on hub: ", err)
	}
	if rep != "ok" {
		log.Fatalf("expected the hub to answer with ok. he did with: %+v", rep)
	}
}

func unregisterOnHub(c *rpc2.Client) {
	var rep string
	host, _ := os.Hostname()
	err := c.Call("unregister", host, &rep)
	if err != nil {
		log.Fatalf("error when unregistring from hub: %s", err)
	}
	c.Close()
}

func handleCheckResultToHub(dc *models.DeamonClient) {
	for {
		var r models.CheckResult = <- dc.ToHub
		log.Println("telling hub about it", r)

		if err := dc.Client.Call("checkResult", r, nil); err != nil {
			log.Println("error sending CheckResult to hub", err)
		}
	}
}

func updateChecksCount(dc *models.DeamonClient, reset bool) {
	if (reset) {
		dc.Deamon.ChecksCount = 0
	} else {
		dc.Deamon.ChecksCount++
	}

	var rep string
	err := dc.Client.Call("updateCheckcount", dc.Deamon, &rep)
	if err != nil {
		log.Println("error updating ChecksCount on hub: ", err)
	}
}


