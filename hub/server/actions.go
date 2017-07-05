package server

import (
	"github.com/cenkalti/rpc2"
	"github.com/mitchellh/mapstructure"
	"github.com/oscp/openshift-monitoring/models"
	"log"
)

func daemonLeave(h *Hub, host string) {
	log.Println("daemon left: ", host)
	delete(h.daemons, host)

	h.toUi <- models.BaseModel{Type: models.DAEMON_LEFT, Message: host}
}

func daemonJoin(h *Hub, d *models.Daemon, c *rpc2.Client) {
	log.Println("new daemon joined:", d)

	h.daemons[d.Hostname] = &models.DaemonClient{Client: c, Daemon: *d}

	if h.currentChecks.IsRunning {
		// Tell the new daemon to join the checks
		if err := c.Call("startChecks", h.currentChecks, nil); err != nil {
			log.Println("error starting checks on newly joined daemon", err)
		}
	}

	h.toUi <- models.BaseModel{Type: models.NEW_DAEMON, Message: d.Hostname}
}

func startChecks(h *Hub, msg interface{}) models.BaseModel {
	checks := getChecksStruct(msg)

	// Save current state & tell daemons
	checks.IsRunning = true
	h.currentChecks = checks
	h.startChecks <- checks

	// Return ok to UI
	return models.BaseModel{Type: models.CURRENT_CHECKS, Message: checks}
}

func stopChecks(h *Hub) models.BaseModel {
	// Save current state & tell daemons
	h.currentChecks.IsRunning = false
	h.stopChecks <- true

	// Return ok to UI
	return models.BaseModel{Type: models.CURRENT_CHECKS, Message: h.currentChecks}
}

func getChecksStruct(msg interface{}) models.Checks {
	var checks models.Checks
	err := mapstructure.Decode(msg, &checks)
	if err != nil {
		log.Println("error decoding checks", err)
	}
	return checks
}
