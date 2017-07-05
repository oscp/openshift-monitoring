package main

import (
	"flag"
	"github.com/oscp/openshift-monitoring/hub/server"
	"log"
	"net/http"
)

var uiAddr = flag.String("UI_ADDR", "localhost:8080", "http service endpoint")
var hubAddr = flag.String("RPC_ADDR", "localhost:2600", "go hub rcp2 address")
var masterApiUrls = flag.String("MASTER_API_URLS", "https://master1:8443,https://master2:8443", "addresses of master api's")
var daemonPublicUrl = flag.String("DAEMON_PUBLIC_URL", "http://daemon.yourroute.com", "external address of the daemon service (route)")
var etcdIps = flag.String("ETCD_IPS", "https://localhost:2379,https://server1:2379", "adresses of etcd servers")
var etcdCertPath = flag.String("ETCD_CERT_PATH", "/etc/etcd/", "Path of alternative etcd certificates")

func main() {
	flag.Parse()
	log.Println("hub waiting for daemons on", *hubAddr)
	log.Println("ui server waiting for websocket on", *uiAddr)
	log.Println("master api urls are", *masterApiUrls)
	log.Println("daemons public url is", *daemonPublicUrl)
	log.Println("etcd ips are", *etcdIps)
	log.Println("etcdCertPath is", *etcdCertPath)

	// Start hub rcp server
	hub := server.NewHub(*hubAddr, *masterApiUrls, *daemonPublicUrl, *etcdIps, *etcdCertPath)
	go hub.Serve()

	// Serve UI & websockets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		server.OnUISocket(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(*uiAddr, nil))
}
