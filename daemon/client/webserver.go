package client

import (
	"github.com/oscp/openshift-monitoring/daemon/client/handlers"
	"log"
	"net/http"
	"os"
)

func RunWebserver(daemonType string) {
	addr := os.Getenv("SERVER_ADDRESS")

	if len(addr) == 0 {
		addr = ":8090"
	}

	log.Println("starting webserver on", addr)

	http.HandleFunc("/fast", handlers.FastHandler)
	http.HandleFunc("/slow", handlers.SlowHandler)
	http.HandleFunc("/checks/minor", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMinorChecks(daemonType, w, r)
	})
	http.HandleFunc("/checks/major", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMajorChecks(daemonType, w, r)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
