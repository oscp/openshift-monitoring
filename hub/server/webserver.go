package server

import (
	"github.com/gorilla/websocket"
	"github.com/oscp/openshift-monitoring/models"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func OnUISocket(h *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("ui joined")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade-error: ", err)
		return
	}

	go handleFromUI(h, c)
	go handleToUI(h, c)
}

func handleToUI(h *Hub, c *websocket.Conn) {
	for {
		var msg models.BaseModel = <-h.toUi

		err := c.WriteJSON(msg)
		if err != nil {
			log.Println("socket to UI was closed, resending message", err)
			h.toUi <- msg
			break
		}
	}
}

func handleFromUI(h *Hub, c *websocket.Conn) {
	for {
		// parse message
		var msg models.BaseModel
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read-error on ws: ", err)
			break
		}

		var res interface{}
		switch msg.Type {
		case models.ALL_DAEMONS:
			res = models.BaseModel{Type: models.ALL_DAEMONS, Message: h.Daemons()}
			break
		case models.START_CHECKS:
			res = startChecks(h, msg.Message)
			break
		case models.STOP_CHECKS:
			res = stopChecks(h)
			break
		case models.CURRENT_CHECKS:
			res = models.BaseModel{Type: models.CURRENT_CHECKS, Message: h.currentChecks}
		}

		err = c.WriteJSON(res)
		if err != nil {
			log.Println("error sending message to UI on websocket: ", err)
		}
	}
}
