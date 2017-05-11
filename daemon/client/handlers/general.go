package handlers

import (
	"net/http"
	"io"
	"time"
	"math/rand"
	"github.com/oscp/openshift-monitoring/models"
	"os"
	"encoding/json"
)

func FastHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world")
}

func SlowHandler(w http.ResponseWriter, r *http.Request) {
	s := random(1, 60000)
	time.Sleep(time.Duration(s) * time.Millisecond)

	io.WriteString(w, "Hello, world")
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func generateResponse(w http.ResponseWriter, responses []models.CheckState) {
	host, _ := os.Hostname()
	r := models.CheckResult{
		Hostname: host,
		Type: "MINOR",
	}

	for _, s := range responses {
		if (!s.State) {
			r.IsOk = false
			r.Message += s.Message
		}
	}

	json, err := json.Marshal(r)
	if (err != nil) {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}