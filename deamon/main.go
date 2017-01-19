package main

import (
	"flag"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/deamon/client"
)

var hubAddr = flag.String("hubAddr", "localhost:2600", "go rcp hub address")
var port = flag.String("port", "2601", "client rcp port")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Register on hub
	client.RegisterOnHub(*hubAddr)

	// Start own server for tasks
	client.DeamonServer(*port)
}
