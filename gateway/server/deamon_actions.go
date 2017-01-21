package server

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"log"
)

func deamonLeave(h *Hub, addr string) {
	log.Println("deamon left: ", addr)
	delete(h.deamons, addr)

	h.toUi <- models.BaseModel{Type: models.TYPE_DEAMON_LEFT, Message: addr}
}

func deamonJoin(h *Hub, addr string, d *models.Deamon) {
	log.Printf("new deamon joined: %+v, addr: %s\n", d, addr)

	h.deamons[addr] = *d

	// tell the ui about it
	h.toUi <- models.BaseModel{Type: models.TYPE_NEW_DEAMON, Message: d}
}
