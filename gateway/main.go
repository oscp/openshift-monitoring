package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/gateway/server"
)

var uiAddr = flag.String("uiAddr", "localhost:8080", "http service endpoint")
var hubAddr = flag.String("hubAddr", "localhost:2600", "go hub rcp address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("ui server waiting for websocket on ", *uiAddr)
	log.Println("hub waiting for deamons on ", *hubAddr)

	// Start hub rcp server
	hub := server.NewHub(*hubAddr)
	go hub.Serve()

	// Start websocket server for ui
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		server.OnUISocket(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(*uiAddr, nil))
}

