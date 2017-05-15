package handlers

import (
	"net/http"
	"github.com/oscp/openshift-monitoring/models"
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
	"os"
	"log"
	"strings"
)

func HandleMajorChecks(daemonType string, w http.ResponseWriter, r *http.Request) {
	responses := []models.CheckState{}
	if (daemonType == "NODE") {
		ok, msg := checks.CheckDockerPool(90)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckDnsNslookupOnKubernetes()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckDnsServiceNode()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	if (daemonType == "MASTER") {
		etcdIps := os.Getenv("ETCD_IPS")
		registryIp := os.Getenv("REGISTRY_SVC_IP")
		routerIps := os.Getenv("ROUTER_IPS")
		if (len(etcdIps) == 0 || len(registryIp) == 0 || len(routerIps) == 0) {
			log.Fatal("env variables 'ETCD_IPS', 'REGISTRY_SVC_IP', 'ROUTER_IPS' must be specified on type 'MASTER'")
		}

		ok, msg := checks.CheckOcGetNodes()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckEtcdHealth(etcdIps, "")
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckRegistryHealth(registryIp)
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		for _, rip := range strings.Split(routerIps, ",") {
			ok, msg = checks.CheckRouterHealth(rip)
			responses = append(responses, models.CheckState{
				State: ok,
				Message: msg,
			})
		}

		ok, msg = checks.CheckMasterApis("https://localhost:8443/api")
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckLoggingRestartsCount()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckDnsNslookupOnKubernetes()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})

		ok, msg = checks.CheckDnsServiceNode()
		responses = append(responses, models.CheckState{
			State: ok,
			Message: msg,
		})
	}

	if (daemonType == "STORAGE") {
		isGlusterServer := os.Getenv("IS_GLUSTER_SERVER")

		if (len(isGlusterServer) > 0) {
			ok, msg := checks.CheckGlusterStatus()
			responses = append(responses, models.CheckState{
				State: ok,
				Message: msg,
			})
		}
	}

	generateResponse(w, responses)
}