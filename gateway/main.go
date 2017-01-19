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

	go server.DeamonHub(*hubAddr)

	http.HandleFunc("/ui", server.OnUISocket)

	log.Fatal(http.ListenAndServe(*uiAddr, nil))
}

