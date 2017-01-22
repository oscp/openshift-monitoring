package client

import (
	"net/http/httptest"
	"net/http"
	"log"
	"github.com/cenkalti/rpc2"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
)

func startJob(client *rpc2.Client, job *models.Job) {
	log.Println("starting new job", job)
}

func stopJob(client *rpc2.Client, msg *string, reply *string) {
	log.Println("stopping job", msg)
}

func checkHttpConnection(url string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	handler(w, req)

	log.Printf("%d - %s", w.Code, w.Body.String())
}