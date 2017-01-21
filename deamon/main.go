package main

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/deamon/client"
	"os"
	"os/signal"
	"log"
	"syscall"
)

func main() {
	// Get config
	hubAddr := os.Getenv("HUB_ADDRESS")
	deamonType := os.Getenv("DEAMON_TYPE")

	// Register on hub
	cl := client.RegisterOnHub(hubAddr, deamonType)

	// Exit gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	func() {
		<-c
		log.Println("got sigterm, unregistring on hub")
		client.UnregisterOnHub(cl)
		os.Exit(1)
	}()
}
