package client

import (
	"net/http"
	"io"
	"time"
	"math/rand"
	"log"
)

func ServeWeb() {
	log.Println("starting webserver on :8090")
	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", slowHandler)
	go log.Fatal(http.ListenAndServe(":8090", nil))
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
	return rand.Intn(max - min) + min
}