package server

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"log"
	"github.com/cenkalti/rpc2"
	"github.com/mitchellh/mapstructure"
)

func deamonLeave(h *Hub, host string) {
	log.Println("deamon left: ", host)
	delete(h.deamons, host)

	h.toUi <- models.BaseModel{Type: models.DEAMON_LEFT, Message: host}
}

func deamonJoin(h *Hub, d *models.Deamon, c *rpc2.Client) {
	log.Println("new deamon joined:", d)

	h.deamons[d.Hostname] = models.DeamonClient{Client:c, Deamon: *d}

	h.toUi <- models.BaseModel{Type: models.NEW_DEAMON, Message: d.Hostname}
}

func startChecks(h *Hub, msg interface{}) models.BaseModel {
	checks := getChecksStruct(msg)

	// Save current state & tell deamons
	checks.IsRunning = true
	h.currentChecks = checks
	h.startChecks <- checks

	// Return ok to UI
	return models.BaseModel{Type: models.START_CHECKS, Message: "ok"}
}

func stopChecks(h *Hub) models.BaseModel {
	// Save current state & tell deamons
	h.currentChecks = models.Checks{IsRunning:false}
	h.stopChecks <- true

	// Return ok to UI
	return models.BaseModel{Type: models.STOP_CHECKS, Message: "ok"}
}

func getChecksStruct(msg interface{}) models.Checks {
	var checks models.Checks
	err := mapstructure.Decode(msg, &checks)
	if err != nil {
		log.Println("error decoding checks", err)
	}
	return checks
}