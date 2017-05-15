package handlers

import (
	"net/http"
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
	"github.com/oscp/openshift-monitoring/models"
	"os"
	"log"
	"strconv"
)

func HandleMinorChecks(daemonType string, w http.ResponseWriter, r *http.Request) {
	responses := []models.CheckState{}
	if (daemonType == "NODE") {
		ok, msg := checks.CheckDockerPool(80)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckHttpService(false)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	if (daemonType == "MASTER") {
		externalSystem := os.Getenv("EXTERNAL_SYSTEM_URL")
		hawcularIp := os.Getenv("HAWCULAR_SVC_IP")
		allowedWithout := os.Getenv("PROJECTS_WITHOUT_LIMITS")
		if (len(externalSystem) == 0 || len(allowedWithout) == 0) {
			log.Fatal("env variables 'EXTERNAL_SYSTEM_URL', 'PROJECTS_WITHOUT_LIMITS', 'HAWCULAR_SVC_IP' must be specified on type 'MASTER'")
		}

		allowedWithoutInt, err := strconv.Atoi(allowedWithout)
		if (err != nil) {
			log.Fatal("allowedWithout seems not to be an integer", allowedWithout)
		}

		ok, msg := checks.CheckExternalSystem(externalSystem)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckHawcularHealth(hawcularIp)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckRouterRestartCount()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckLimitsAndQuotas(allowedWithoutInt)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckHttpService(false)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	if (daemonType == "STORAGE") {
		ok, msg := checks.CheckOpenFileCount()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	ok, msg := checks.CheckNtpd()
	responses = append(responses, models.CheckState{
		State: ok,
		Message: msg,
	})

	generateResponse(w, responses)
}
