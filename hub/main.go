package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/hub/server"
)

var uiAddr = flag.String("UI_ADDR", "localhost:8080", "http service endpoint")
var hubAddr = flag.String("RPC_ADDR", "localhost:2600", "go hub rcp2 address")
var masterApiUrls = flag.String("MASTER_API_URLS", "http://www.google.ch,http://www.heise.de", "addresses of master api's")

func main() {
	flag.Parse()
	log.Println("hub waiting for deamons on ", *hubAddr)
	log.Println("ui server waiting for websocket on ", *uiAddr)
	log.Println("master api urls are ", *masterApiUrls)

	// Start hub rcp server
	hub := server.NewHub(*hubAddr, *masterApiUrls)
	go hub.Serve()

	// Start websocket server for ui
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		server.OnUISocket(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(*uiAddr, nil))
}

