package server

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GuiRequest struct {
	Type string
	Message string
}

func OnUISocket(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection from UI")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade-error: ", err)
		return
	}
	defer c.Close()

	for {
		// parse message
		var msg GuiRequest
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read-error on ws: ", err)
			break
		}

		log.Println("new message from client: ", msg.Type, msg.Message)

		var res interface{}
		switch msg.Type {
		case "getKeywords":
			break
		case "getResults":
			break
		case "newKeyword":
			break
		case "deleteKeyword" :
			break
		}

		err = c.WriteJSON(res)
		if err != nil {
			log.Println("write-error: ", err)
			break
		}
	}
}