package main

import (
	"flag"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/deamon/client"
	"os"
	"os/signal"
	"syscall"
)

var hubAddr = flag.String("hubAddr", "localhost:2600", "go rcp hub address")
var deamonType = flag.String("deamonType", "master", "type of deamon: master,node,pod")
var port = flag.Int("port", 2601, "client rcp port")

func main() {
	flag.Parse()
	log.SetFlags(0)

	// Register on hub
	cl := client.RegisterOnHub(*hubAddr, *deamonType, *port)

	// Exit gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.UnregisterOnHub(cl)
		os.Exit(1)
	}()

	// Start own server for tasks
	client.DeamonServer(*port)
}
