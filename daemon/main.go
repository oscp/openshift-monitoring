package main

import (
	"github.com/oscp/openshift-monitoring/daemon/client"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	daemonType := os.Getenv("DAEMON_TYPE")

	if len(daemonType) == 0 {
		log.Fatal("env variable 'DAEMON_TYPE' must be specified")
	}

	// Webserver for /slow /fast checks
	go client.RunWebserver()

	hubAddr := os.Getenv("HUB_ADDRESS")
	namespace := os.Getenv("POD_NAMESPACE")

	if daemonType == "POD" && len(namespace) == 0 {
		log.Fatal("if type is 'POD' env variable 'POD_NAMESPACE' must be specified")
	}

	// If no hub address is provided, start only the webserver
	if len(hubAddr) > 0 {
		// Register on hub
		cl := client.StartDaemon(hubAddr, daemonType, namespace)

		// Exit gracefully
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		func() {
			<-c
			log.Println("got sigterm, unregistring on hub")
			client.StopDaemon(cl)
			os.Exit(1)
		}()
	} else {
		// Sleep 4 ever
		select {}
	}
}
