package client

import (
	"net/http"
	"io"
)

func ServeWeb() {
	http.HandleFunc("/hello", helloHandler)
	go http.ListenAndServe(":8090", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world")
}