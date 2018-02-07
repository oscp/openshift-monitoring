package client

import (
	"log"
	"net/http"
	"io"
	"time"
	"math/rand"
)

func RunWebserver() {
	addr := ":8090"
	log.Println("starting webserver on", addr)

	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", slowHandler)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	s := random(1, 60000)
	time.Sleep(time.Duration(s) * time.Millisecond)

	io.WriteString(w, "Hello, world")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}