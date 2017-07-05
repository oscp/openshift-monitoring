package handlers

import (
	"encoding/json"
	"github.com/oscp/openshift-monitoring/models"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
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
	return rand.Intn(max-min) + min
}

func generateResponse(w http.ResponseWriter, errors []string) {
	host, _ := os.Hostname()
	r := models.CheckResult{
		Hostname: host,
		Type:     "OSE_CHECKS",
		IsOk:     true,
	}

	for _, s := range errors {
		r.IsOk = false
		r.Message += " | " + s
	}

	json, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Error while generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
