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
	namespace := os.Getenv("POD_NAMESPACE")

	if (len(hubAddr) == 0 || len(deamonType) == 0 || len(namespace) == 0) {
		log.Fatal("env variables 'HUB_ADDRESS', 'DEAMON_TYPE', 'POD_NAMESPACE' must be specified")
	}

	// Register on hub
	cl := client.StartDeamon(hubAddr, deamonType, namespace)

	// Create webserver for checks
	go client.ServeWeb()


	// Exit gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	func() {
		<-c
		log.Println("got sigterm, unregistring on hub")
		client.StopDeamon(cl)
		os.Exit(1)
	}()
}
