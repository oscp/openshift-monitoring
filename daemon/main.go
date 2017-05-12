package main

import (
	"os"
	"log"
	"github.com/oscp/openshift-monitoring/daemon/client"
	"syscall"
	"os/signal"
)

func main() {
	daemonType := os.Getenv("DAEMON_TYPE")
	withHub := os.Getenv("WITH_HUB")

	if (len(withHub) == 0 || len(daemonType) == 0) {
		log.Fatal("env variables 'WITH_HUB' and 'DAEMON_TYPE' must be specified")
	}

	// Communication with the hub is optional
	if (withHub == "true") {
		// Webserver for /slow /fast checks
		go client.RunWebserver(daemonType)

		hubAddr := os.Getenv("HUB_ADDRESS")
		namespace := os.Getenv("POD_NAMESPACE")

		if (len(hubAddr) == 0) {
			log.Fatal("env variable 'HUB_ADDRESS' must be specified")
		}

		if (daemonType == "POD" && len(namespace) == 0) {
			log.Fatal("if type is 'POD' env variable 'POD_NAMESPACE' must be specified")
		}

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
		// Just run the webserver for external monitoring system
		client.RunWebserver(daemonType)
	}
}
