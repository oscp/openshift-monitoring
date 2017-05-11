package handlers

import (
	"net/http"
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
	"github.com/oscp/openshift-monitoring/models"
)

func HandleMinorChecks(daemonType string, w http.ResponseWriter, r *http.Request) {
	responses := []models.CheckState{}
	if (daemonType == "NODE") {
		ok, msg := checks.CheckDockerPool(80)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	generateResponse(w, responses)
}
