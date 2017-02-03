package client

import (
	"github.com/oscp/openshift-monitoring/models"
	"log"
	"os"
	"github.com/cenkalti/rpc2"
)

func registerOnHub(h string, dc *models.DaemonClient) {
	log.Println("registring on the hub:", h)
	var rep string
	err := dc.Client.Call("register", dc.Daemon, &rep)
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

func handleCheckStarted(dc *models.DaemonClient) {
	dc.Daemon.StartedChecks++
	updateDaemonOnHub(dc)
}

func handleCheckFinished(dc *models.DaemonClient, ok bool) {
	if (ok) {
		dc.Daemon.SuccessfulChecks++
	} else {
		dc.Daemon.FailedChecks++
	}
	updateDaemonOnHub(dc)
}

func handleChecksStopped(dc *models.DaemonClient) {
	log.Println("stopped checks")

	dc.Daemon.StartedChecks = 0;
	dc.Daemon.FailedChecks = 0;
	dc.Daemon.SuccessfulChecks = 0;
	updateDaemonOnHub(dc)
}

func updateDaemonOnHub(dc *models.DaemonClient) {
	var rep string
	err := dc.Client.Call("updateCheckcount", dc.Daemon, &rep)
	if err != nil {
		log.Println("error updating Checkcounts on hub: ", err)
	}
}

func handleCheckResultToHub(dc *models.DaemonClient) {
	for {
		var r models.CheckResult = <- dc.ToHub
		r.Hostname = dc.Daemon.Hostname

		if err := dc.Client.Call("checkResult", r, nil); err != nil {
			log.Println("error sending CheckResult to hub", err)
		}
	}
}


