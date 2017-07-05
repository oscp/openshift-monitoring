package client

import (
	"github.com/cenkalti/rpc2"
	"github.com/oscp/openshift-monitoring/models"
	"log"
	"os"
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

func HandleCheckStarted(dc *models.DaemonClient) {
	dc.Daemon.StartedChecks++
	updateDaemonOnHub(dc)
}

func HandleCheckFinished(dc *models.DaemonClient, err error, t string) {
	// Update check counts
	if err == nil {
		dc.ToHub <- models.CheckResult{Type: t, IsOk: true, Message: ""}
		dc.Daemon.SuccessfulChecks++
	} else {
		dc.ToHub <- models.CheckResult{Type: t, IsOk: false, Message: err.Error()}
		dc.Daemon.FailedChecks++
	}
	updateDaemonOnHub(dc)
}

func HandleChecksStopped(dc *models.DaemonClient) {
	log.Println("stopped checks")
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
		var r models.CheckResult = <-dc.ToHub
		r.Hostname = dc.Daemon.Hostname

		if err := dc.Client.Call("checkResult", r, nil); err != nil {
			log.Println("error sending CheckResult to hub", err)
		}
	}
}
