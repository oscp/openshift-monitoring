package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/oscp/openshift-monitoring/daemon/client/checks"
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

	if daemonType == "MASTER" || daemonType == "NODE" {
		certPaths := os.Getenv("CHECK_CERTIFICATE_PATHS")
		kubePaths := os.Getenv("CHECK_CERTIFICATE_KUBE_PATHS")

		if len(certPaths) == 0 || len(kubePaths) == 0 {
			log.Fatal("env variables 'CHECK_CERTIFICATE_PATHS', 'CHECK_CERTIFICATE_KUBE_PATHS' must be specified")
		}

		if err := checks.CheckFileSslCertificates(strings.Split(certPaths, ","), 80); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckKubeSslCertificates(strings.Split(kubePaths, ","), 80); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckBondNetworkInterface(); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "MASTER" {
		externalSystem := os.Getenv("EXTERNAL_SYSTEM_URL")
		hawcularIp := os.Getenv("HAWCULAR_SVC_IP")
		allowedWithoutLimits := os.Getenv("PROJECTS_WITHOUT_LIMITS")
		allowedWithoutQuota := os.Getenv("PROJECTS_WITHOUT_QUOTA")
		certUrls := os.Getenv("CHECK_CERTIFICATE_URLS")

		if len(externalSystem) == 0 || len(allowedWithoutLimits) == 0 || len(allowedWithoutQuota) == 0 || len(certUrls) == 0 {
			log.Fatal("env variables 'EXTERNAL_SYSTEM_URL', 'PROJECTS_WITHOUT_LIMITS', 'PROJECTS_WITHOUT_QUOTA', 'CHECK_CERTIFICATE_URLS' must be specified on type 'MASTER'")
		}

		allowedWithoutLimitsInt, err := strconv.Atoi(allowedWithoutLimits)
		if err != nil {
			log.Fatal("allowedWithoutLimits seems not to be an integer", allowedWithoutLimits)
		}
		allowedWithoutQuotaInt, err := strconv.Atoi(allowedWithoutQuota)
		if err != nil {
			log.Fatal("allowedWithoutLimits seems not to be an integer", allowedWithoutQuota)
		}

		// boolean false means exclude buildnodes
		// boolean true means only buildnodes
		if err := checks.CheckOcGetNodes(true); err != nil {
			errors = append(errors, err.Error())
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

		if err := checks.CheckLimitsAndQuota(allowedWithoutLimitsInt, allowedWithoutQuotaInt); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckHttpService(false); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckLoggingRestartsCount(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckUrlSslCertificates(strings.Split(certUrls, ","), 80); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "STORAGE" {
		if err := checks.CheckOpenFileCount(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckMountPointSizes(85); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckLVPoolSizes(80); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckVGSizes(10); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if err := checks.CheckChrony(); err != nil {
		errors = append(errors, err.Error())
	}

	generateResponse(w, errors)
}
