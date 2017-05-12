package checks

import (
	"strings"
	"bytes"
	"log"
	"os/exec"
	"strconv"
)

func CheckMasterApis(urls string) (bool, string) {
	urlArr := strings.Split(urls, ",")

	oneApiOk := false
	var msg string
	for _, u := range urlArr {
		if (checkHttp(u)) {
			oneApiOk = true
		} else {
			msg += u + " is not reachable. ";
		}
	}

	return oneApiOk, msg
}

func CheckOcGetNodes() (bool, string) {
	var msg string
	out, err := exec.Command("bash", "-c", "oc get nodes --show-labels | grep -v monitoring=false").Output()
	if err != nil {
		msg = "Could not parse oc get nodes output: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	isOk := !strings.Contains(string(out), "NotReady")
	if (!isOk) {
		msg = "Some node is not ready! 'oc get nodes' output contained NotReady"
	}
	return isOk, msg
}

func CheckDnsNslookupOnKubernetes() (bool, string) {
	var msg string

	cmd := exec.Command("nslookup", daemonDNSEndpoint, kubernetesIP)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("error with nslookup: ", err)
		msg = "DNS resolution via nslookup & kubernetes failed." + err.Error()
		return false, msg
	}

	stdOut := out.String()

	isOk := false
	if (strings.Contains(stdOut, "Server") && strings.Count(stdOut, "Address") >= 2 && strings.Contains(stdOut, "Name")) {
		isOk = true
	} else {
		msg += "NsLookup had wrong output"
	}

	return isOk, msg
}

func CheckDnsServiceNode() (bool, string) {
	isOk := false
	var msg string

	ips := getIpsForName(daemonDNSServiceA)

	if (ips == nil) {
		msg = "Failed to lookup ip on node (dnsmasq) for name " + daemonDNSServiceA
	} else {
		isOk = true
	}

	return isOk, msg
}

func CheckDnsInPod() (bool, string) {
	isOk := false
	var msg string

	ips := getIpsForName(daemonDNSPod)

	if (ips == nil) {
		msg = "Failed to lookup ip in pod for name " + daemonDNSPod
	} else {
		isOk = true
	}

	return isOk, msg
}

func CheckPodHttpAtoB() (bool, string) {
	// This should fail as we do not have access to this project
	isOk := !checkHttp("http://" + daemonDNSServiceB + ":8090/hello")

	var msg string
	if (!isOk) {
		msg = "Pod A could access pod b. This should not be allowed!"
	}

	return isOk, msg
}

func CheckPodHttpAtoC(slow bool) (bool, string) {
	var msg string
	isOk := checkHttp("http://" + daemonDNSServiceC + ":8090/" + getEndpoint(slow))

	if (!isOk) {
		msg = "Pod A could access pod c. Route/Router problem?"
	}

	return isOk, msg
}

func CheckHttpService(slow bool) (bool, string) {
	var msg string

	isOkA := checkHttp("http://" + daemonDNSServiceA + ":8090/" + getEndpoint(slow))
	isOkB := checkHttp("http://" + daemonDNSServiceB + ":8090/" + getEndpoint(slow))
	isOkC := checkHttp("http://" + daemonDNSServiceC + ":8090/" + getEndpoint(slow))

	isOk := true
	if (!isOkA || !isOkB || !isOkC) {
		msg = "Could not reach one of the services (a/b/c)"
		isOk = false
	}

	return isOk, msg
}

func CheckHttpHaProxy(publicUrl string, slow bool) (bool, string) {
	var msg string

	isOk := checkHttp(publicUrl + ":80/" + getEndpoint(slow))

	if (!isOk) {
		msg = "Could not access pods via haproxy. Route/Router problem?"
	}

	return isOk, msg
}

func CheckRegistryHealth(ip string) (bool, string) {
	var msg string
	isOk := checkHttp("http://" + ip + ":5000/healthz")

	if (!isOk) {
		msg = "Registry health check failed"
	}

	return isOk, msg
}

func CheckHawcularHealth(ip string) (bool, string) {
	var msg string
	isOk := checkHttp("https://" + ip + ":443")

	if (!isOk) {
		msg = "Hawcular health check failed"
	}

	return isOk, msg
}

func CheckRouterHealth(ip string) (bool, string) {
	var msg string
	isOk := checkHttp("http://" + ip + ":1936/healthz")

	if (!isOk) {
		msg = "Router health check failed:" + ip
	}

	return isOk, msg
}

func CheckLoggingRestartsCount() (bool, string) {
	var msg string
	out, err := exec.Command("bash", "-c", "oc get pods -n logging -o wide | tr -s ' ' | cut -d ' ' -f 4").Output()
	if err != nil {
		msg = "Could not parse logging container restart count: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	isOk := true
	for _, l := range strings.Split(string(out), "\n") {
		if (!strings.HasPrefix(l, "RESTARTS") && len(strings.TrimSpace(l)) > 0) {
			cnt, _ := strconv.Atoi(l)
			if (cnt > 2) {
				msg = "A logging-container has restart count > 2"
				isOk = false
			}
		}
	}

	return isOk, msg
}

func CheckRouterRestartCount() (bool, string) {
	var msg string
	out, err := exec.Command("bash", "-c", "oc get po -n default | grep router | grep -v deploy | tr -s ' ' | cut -d ' ' -f 4").Output()
	if err != nil {
		msg = "Could not parse router restart count: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	isOk := true
	for _, l := range strings.Split(string(out), "\n") {
		if (!strings.HasPrefix(l, "RESTARTS") && len(strings.TrimSpace(l)) > 0) {
			cnt, _ := strconv.Atoi(l)
			if (cnt > 5) {
				msg = "A Router has restart count > 5"
				isOk = false
			}
		}
	}

	return isOk, msg
}

func CheckEtcdHealth(etcdIps string, etcdCertPath string) (bool, string) {
	var msg string
	isOk := true

	if (len(etcdCertPath) > 0) {
		// Check etcd with custom certs path
		isOk = checkEtcdHealthWithCertPath(&msg, etcdCertPath, etcdIps)

		if (!isOk) {
			log.Println("etcd health check with custom cert path failed, trying with default")

			// Check etcd with default certs path
			isOk = checkEtcdHealthWithCertPath(&msg, "/etc/etcd/", etcdIps);
		}
	} else {
		// Check etcd with default certs path
		isOk = checkEtcdHealthWithCertPath(&msg, "/etc/etcd/", etcdIps);
	}
	return isOk, msg
}

func checkEtcdHealthWithCertPath(msg *string, certPath string, etcdIps string) bool {
	cmd := exec.Command("etcdctl", "--peers", etcdIps, "--ca-file", certPath + "ca.crt",
		"--key-file", certPath + "peer.key", "--cert-file", certPath + "peer.crt", "cluster-health")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("error while running etcd health check", err)
		*msg = "etcd health check failed: " + err.Error()
		return false
	}

	stdOut := out.String()
	if (!strings.Contains(stdOut, "cluster is healthy")) {
		*msg += "Etcd health check was 'cluster unhealthy'"
		return false
	}

	return true
}

func CheckLimitsAndQuotas(allowedWithout int) (bool, string) {
	var msg string

	// Count projects
	projectCount, err := exec.Command("bash", "-c", "oc get projects  | wc -l").Output()
	if err != nil {
		msg = "Could not parse project count" + err.Error()
		log.Println(msg)
		return false, msg
	}

	// Count limits
	limitCount, err := exec.Command("bash", "-c", "oc get limits --all-namespaces | wc -l").Output()
	if err != nil {
		msg = "Could not parse limit count" + err.Error()
		log.Println(msg)
		return false, msg
	}

	// Count quotas
	quotaCount, err := exec.Command("bash", "-c", "oc get quota --all-namespaces | wc -l").Output()
	if err != nil {
		msg = "Could not parse quota count" + err.Error()
		log.Println(msg)
		return false, msg
	}

	// Parse them
	pCount, err := strconv.Atoi(strings.TrimSpace(string(projectCount)))
	lCount, _ := strconv.Atoi(strings.TrimSpace(string(limitCount)))
	qCount, _ := strconv.Atoi(strings.TrimSpace(string(quotaCount)))

	log.Println("Parsed values (projects,limits,quotas)", pCount, lCount, qCount)

	if (pCount - allowedWithout != lCount) {
		return false, "There are some projects without limits"
	}
	if (pCount - allowedWithout != qCount) {
		return false, "There are some projects without quotas"
	}

	return true, ""
}
