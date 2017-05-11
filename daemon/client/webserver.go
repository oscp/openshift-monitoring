package client

import (
	"net/http"
	"log"
	"os"
	"github.com/oscp/openshift-monitoring/daemon/client/handlers"
)

func RunWebserver(daemonType string) {
	addr := os.Getenv("SERVER_ADDRESS")

	if (len(addr) == 0) {
		addr = ":8090"
	}

	log.Println("starting webserver on", addr)

	http.HandleFunc("/fast", handlers.FastHandler)
	http.HandleFunc("/slow", handlers.SlowHandler)
	http.HandleFunc("/checks/minor", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMinorChecks(daemonType, w, r)
	})

	go log.Fatal(http.ListenAndServe(addr, nil))
}
