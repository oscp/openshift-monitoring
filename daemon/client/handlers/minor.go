package handlers

import (
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
	"log"
	"net/http"
	"os"
	"strconv"
)

func HandleMinorChecks(daemonType string, w http.ResponseWriter, r *http.Request) {
	errors := []string{}
	if daemonType == "NODE" {
		if err := checks.CheckDockerPool(80); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckHttpService(false); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "MASTER" {
		externalSystem := os.Getenv("EXTERNAL_SYSTEM_URL")
		hawcularIp := os.Getenv("HAWCULAR_SVC_IP")
		allowedWithout := os.Getenv("PROJECTS_WITHOUT_LIMITS")
		if len(externalSystem) == 0 || len(allowedWithout) == 0 {
			log.Fatal("env variables 'EXTERNAL_SYSTEM_URL', 'PROJECTS_WITHOUT_LIMITS', 'HAWCULAR_SVC_IP' must be specified on type 'MASTER'")
		}

		allowedWithoutInt, err := strconv.Atoi(allowedWithout)
		if err != nil {
			log.Fatal("allowedWithout seems not to be an integer", allowedWithout)
		}

		if err := checks.CheckExternalSystem(externalSystem); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckHawcularHealth(hawcularIp); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckRouterRestartCount(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckLimitsAndQuotas(allowedWithoutInt); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckHttpService(false); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckLoggingRestartsCount(); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "STORAGE" {
		if err := checks.CheckOpenFileCount(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckLVPoolSizes(80); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckVGSizes(10); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if err := checks.CheckNtpd(); err != nil {
		errors = append(errors, err.Error())
	}

	generateResponse(w, errors)
}
