package server

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"log"
	"github.com/cenkalti/rpc2"
)

func deamonLeave(h *Hub, host string) {
	log.Println("deamon left: ", host)
	delete(h.deamons, host)

	h.toUi <- models.BaseModel{Type: models.TYPE_DEAMON_LEFT, Message: host}
}

func deamonJoin(h *Hub, d *models.Deamon, c *rpc2.Client) {
	log.Println("new deamon joined:", d)

	h.deamons[d.Hostname] = models.DeamonClient{Client:c, Deamon: *d}

	// tell the ui about it
	h.toUi <- models.BaseModel{Type: models.TYPE_NEW_DEAMON, Message: d.Hostname}
}
