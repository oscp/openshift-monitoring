package client

import (
	"net/http/httptest"
	"net/http"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
)

func startChecks(dc *models.DeamonClient, checks *models.Checks) {
	log.Println("starting new checks", checks)

	dc.Deamon.ChecksCount++
	updateChecksCount(dc)
}

func stopChecks(dc *models.DeamonClient) {
	log.Println("stopping all checks")

	if (dc.Deamon.ChecksCount > 0) {
		dc.Deamon.ChecksCount--
	}
	updateChecksCount(dc)
}

func updateChecksCount(dc *models.DeamonClient) {
	var rep string
	err := dc.Client.Call("updateCheckcount", dc.Deamon, &rep)
	if err != nil {
		log.Fatal("error updating ChecksCount on hub: ", err)
	}
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