package handlers

import (
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
	"log"
	"net/http"
	"os"
	"strings"
)

func HandleMajorChecks(daemonType string, w http.ResponseWriter, r *http.Request) {
	errors := []string{}
	if daemonType == "NODE" {
		if err := checks.CheckDockerPool(90); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckDnsNslookupOnKubernetes(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckDnsServiceNode(); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "MASTER" {
		etcdIps := os.Getenv("ETCD_IPS")
		registryIp := os.Getenv("REGISTRY_SVC_IP")
		routerIps := os.Getenv("ROUTER_IPS")
		if len(etcdIps) == 0 || len(registryIp) == 0 || len(routerIps) == 0 {
			log.Fatal("env variables 'ETCD_IPS', 'REGISTRY_SVC_IP', 'ROUTER_IPS' must be specified on type 'MASTER'")
		}

		if err := checks.CheckOcGetNodes(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckEtcdHealth(etcdIps, ""); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckRegistryHealth(registryIp); err != nil {
			errors = append(errors, err.Error())
		}

		for _, rip := range strings.Split(routerIps, ",") {
			if err := checks.CheckRouterHealth(rip); err != nil {
				errors = append(errors, err.Error())
			}
		}

		if err := checks.CheckMasterApis("https://localhost:8443/api"); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckDnsNslookupOnKubernetes(); err != nil {
			errors = append(errors, err.Error())
		}

		if err := checks.CheckDnsServiceNode(); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if daemonType == "STORAGE" {
		isGlusterServer := os.Getenv("IS_GLUSTER_SERVER")

		if isGlusterServer == "true" {
			if err := checks.CheckGlusterStatus(); err != nil {
				errors = append(errors, err.Error())
			}

			if err := checks.CheckLVPoolSizes(90); err != nil {
				errors = append(errors, err.Error())
			}

			if err := checks.CheckVGSizes(5); err != nil {
				errors = append(errors, err.Error())
			}
		}
	}

	generateResponse(w, errors)
}
