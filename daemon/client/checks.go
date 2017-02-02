package client

import (
	"net/http"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"time"
	"strings"
	"os/exec"
	"bytes"
	"net"
	"crypto/tls"
)

const (
	daemonDNSEndpoint = "daemon.ose-mon-a.endpoints.cluster.local"
	daemonDNSServiceA = "daemon.ose-mon-a.svc.cluster.local"
	daemonDNSServiceB = "daemon.ose-mon-b.svc.cluster.local"
	daemonDNSServiceC = "daemon.ose-mon-c.svc.cluster.local"
	daemonDNSPod = "daemon"
	kubernetesIP = "172.30.0.1"
)

func startChecks(dc *models.DaemonClient, checks *models.Checks) {
	tickExt := time.Tick(time.Duration(checks.CheckInterval) * time.Millisecond)
	tickInt := time.Tick(3 * time.Second)

	log.Println("starting checks")

	go func() {
		for {
			select {
			case <-dc.Quit:
				handleChecksStopped(dc)
				return
			case <-tickInt:
				if (checks.MasterApiCheck) {
					go checkMasterApis(dc, checks.MasterApiUrls)
				}
				if (checks.EtcdCheck && dc.Daemon.IsMaster()) {
					go checkEtcdHealth(dc, checks.EtcdIps)
				}
			case <-tickExt:
				if (checks.DnsCheck) {
					go checkDnsNslookupOnKubernetes(dc)

					if (dc.Daemon.IsNode() || dc.Daemon.IsMaster()) {
						go checkDnsServiceNode(dc)
					}

					if (dc.Daemon.IsPod()) {
						go checkDnsInPod(dc)
					}
				}

				if (checks.HttpChecks) {
					if (dc.Daemon.IsPod() && strings.HasSuffix(dc.Daemon.Namespace, "a")) {
						go checkPodHttpAtoB(dc)
						go checkPodHttpAtoC(dc)
					}

					if (dc.Daemon.IsNode() || dc.Daemon.IsMaster()) {
						go checkHttpService(dc)
					}

					go checkHttpHaProxy(dc, checks.DaemonPublicUrl)
				}
			}
		}
	}()
}

func stopChecks(dc *models.DaemonClient) {
	dc.Quit <- true
}

func checkDnsNslookupOnKubernetes(dc *models.DaemonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	cmd := exec.Command("nslookup", daemonDNSEndpoint, kubernetesIP)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		isOk = false
		log.Println("error with nslookup: ", err)
		msg = "DNS resolution via nslookup & kubernetes failed." + err.Error()
	}

	stdOut := out.String()

	if (strings.Contains(stdOut, "Server") && strings.Count(stdOut, "Address") >= 2 && strings.Contains(stdOut, "Name")) {
		isOk = true
	} else {
		msg += "NsLookup had wrong output"
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.DNS_NSLOOKUP_KUBERNETES, IsOk: isOk, Message: msg}
}

func checkDnsServiceNode(dc *models.DaemonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	ips := getIpsForName(daemonDNSServiceA)

	if (ips == nil) {
		msg = "Failed to lookup ip on node (dnsmasq) for name " + daemonDNSServiceA
	} else {
		isOk = true
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.DNS_SERVICE_NODE, IsOk: isOk, Message: msg}
}

func checkDnsInPod(dc *models.DaemonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	ips := getIpsForName(daemonDNSPod)

	if (ips == nil) {
		msg = "Failed to lookup ip in pod for name " + daemonDNSPod
	} else {
		isOk = true
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.DNS_SERVICE_POD, IsOk: isOk, Message: msg}
}

func getIpsForName(n string) []net.IP {
	ips, err := net.LookupIP(n)
	if (err != nil) {
		log.Println("failed to lookup ip for name ", n)
		return nil
	}
	return ips
}

func checkMasterApis(dc *models.DaemonClient, urls string) {
	handleCheckStarted(dc)
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

	handleCheckFinished(dc, oneApiOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.MASTER_API_CHECK, IsOk: oneApiOk, Message: msg}
}

func checkHttp(toCall string) bool {
	if (strings.HasPrefix(toCall, "https")) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		_, err := client.Get(toCall)
		if (err != nil) {
			log.Println("error in http check: ", err)
		}
		return err == nil
	} else {
		_, err := http.Get(toCall)
		if (err != nil) {
			log.Println("error in http check: ", err)
		}
		return err == nil
	}
}

func checkPodHttpAtoB(dc *models.DaemonClient) {
	// This should fail as we do not have access to this project
	handleCheckStarted(dc)
	var msg string

	isOk := !checkHttp("http://" + daemonDNSServiceB + ":8090/hello")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_POD_SERVICE_A_B, IsOk: isOk, Message: msg}
}

func checkPodHttpAtoC(dc *models.DaemonClient) {
	// This should work as we joined this projects
	handleCheckStarted(dc)
	var msg string

	isOk := checkHttp("http://" + daemonDNSServiceC + ":8090/hello")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_POD_SERVICE_A_C, IsOk: isOk, Message: msg}
}

func checkHttpService(dc *models.DaemonClient) {
	handleCheckStarted(dc)
	var msg string

	isOkA := checkHttp("http://" + daemonDNSServiceA + ":8090/hello")
	isOkB := checkHttp("http://" + daemonDNSServiceB + ":8090/hello")
	isOkC := checkHttp("http://" + daemonDNSServiceC + ":8090/hello")

	isOk := true
	if (!isOkA || !isOkB || !isOkC) {
		msg = "Could not reach one of the services (a/b/c)"
		isOk = false
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_SERVICE_ABC, IsOk: isOk, Message: msg}
}

func checkHttpHaProxy(dc *models.DaemonClient, publicUrl string) {
	handleCheckStarted(dc)
	var msg string

	isOk := checkHttp(publicUrl + ":80/hello")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_HAPROXY, IsOk: isOk, Message: msg}
}

func checkEtcdHealth(dc *models.DaemonClient, etcdIps string) {
	handleCheckStarted(dc)
	var msg string
	isOk := true

	cmd := exec.Command("etcdctl", "--peers", etcdIps, "--ca-file", "/etc/etcd/ca.crt",
		"--key-file", "/etc/etcd/peer.key", "--cert-file", "/etc/etcd/peer.crt", "cluster-health")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		isOk = false
		log.Println("error while running etcd health check", err)
		msg = "etcd health check failed: " + err.Error()
	}

	stdOut := out.String()
	if (!strings.Contains(stdOut, "cluster is healthy")) {
		isOk = false
		msg += "Etcd health check was 'cluster unhealthy'"
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.ETCD_HEALTH, IsOk: isOk, Message: msg}
}