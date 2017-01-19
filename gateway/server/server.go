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

func OnUISocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("new connection from UI")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade-error: ", err)
		return
	}
	defer c.Close()

	// handle message from the ui
	go handleFromUI(c)

	// handle messages to the ui
	go handleToUI(hub)
}

func handleToUI(hub *Hub) {
	for {
		log.Println("waiting for data to send to the UI")

		var msg models.BaseModel = <-hub.toUi
		log.Println("should send this to client: ", msg.Type.Name, msg.Message)
	}
}

func handleFromUI(c *websocket.Conn) {
	for {
		log.Println("UI joined. Waiting for messages from UI")

		// parse message
		var msg models.BaseModel
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read-error on ws: ", err)
			break
		}

		log.Println("new message from client: ", msg.Type, msg.Message)

		var res interface{}
		//switch msg.Type {
		//case "getDeamons":
		//	break
		//case "getResults":
		//	break
		//case "newKeyword":
		//	break
		//case "deleteKeyword" :
		//	break
		//}

		err = c.WriteJSON(res)
		if err != nil {
			log.Println("write-error: ", err)
			break
		}
	}
}