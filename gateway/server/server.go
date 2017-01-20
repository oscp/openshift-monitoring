package server

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func OnUISocket(h *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("ui joined by websockets")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade-error: ", err)
		return
	}

	// handle message from the ui
	go handleFromUI(h, c)

	// handle messages to the ui
	go handleToUI(h, c)
}

func handleToUI(h *Hub, c *websocket.Conn) {
	for {
		var msg models.BaseModel = <-h.toUi

		err := c.WriteJSON(msg)
		if err != nil {
			log.Println("error sending message to UI on websocket: ", err)
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
		case models.TYPE_ALL_DEAMONS:
			res = models.BaseModel{ Type: models.TYPE_ALL_DEAMONS, Message: h.deamons}
			break
		}

		err = c.WriteJSON(res)
		if err != nil {
			log.Println("error sending message to UI on websocket: ", err)
		}
	}
}