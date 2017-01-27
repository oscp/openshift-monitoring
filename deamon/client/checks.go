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
	deamonDNSEndpoint = "deamon.ose-mon-a.endpoints.cluster.local"
	deamonDNSServiceA = "deamon.ose-mon-a.svc.cluster.local"
	deamonDNSServiceB = "deamon.ose-mon-b.svc.cluster.local"
	deamonDNSServiceC = "deamon.ose-mon-c.svc.cluster.local"
	deamonDNSPod = "deamon"
	kubernetesIP = "172.30.0.1"
)

func startChecks(dc *models.DeamonClient, checks *models.Checks) {
	tickExt := time.Tick(time.Duration(checks.CheckInterval) * time.Second)
	tickInt := time.Tick(5 * time.Second)

	log.Println("starting checks")

	go func() {
		for {
			select {
			case <-dc.Quit:
				log.Println("stopped checks")
				return
			case <-tickInt:
				if (checks.MasterApiCheck) {
					go checkMasterApis(dc, checks.MasterApiUrls)
				}
			case <-tickExt:
				if (checks.DnsCheck) {
					go checkDnsNslookupOnKubernetes(dc)

					if (dc.Deamon.IsNode()) {
						go checkDnsServiceNode(dc)
					}

					if (dc.Deamon.IsPod()) {
						go checkDnsInPod(dc)
					}
				}

				if (checks.HttpChecks) {
					if (dc.Deamon.IsNode() || (dc.Deamon.IsPod() && strings.HasSuffix(dc.Deamon.Namespace, "a"))) {
						go checkPodHttpAtoB(dc)
						go checkPodHttpAtoC(dc)
					}

					go checkHttpHaProxy(dc, checks.DeamonPublicUrl)
				}
			}
		}
	}()
}

func stopChecks(dc *models.DeamonClient) {
	dc.Quit <- true
}

func checkDnsNslookupOnKubernetes(dc *models.DeamonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	cmd := exec.Command("nslookup", deamonDNSEndpoint, kubernetesIP)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		isOk = false
		log.Println("error with nslookup: ", err)
		msg = "DNS resolution via nslookup & kubernetes failed."
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

func checkDnsServiceNode(dc *models.DeamonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	ips := getIpsForName(deamonDNSServiceA)

	if (ips == nil) {
		isOk = false
		msg = "Failed to lookup ip on node (dnsmasq) for name " + deamonDNSServiceA
	}

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.DNS_SERVICE_NODE, IsOk: isOk, Message: msg}
}

func checkDnsInPod(dc *models.DeamonClient) {
	handleCheckStarted(dc)
	isOk := false
	var msg string

	ips := getIpsForName(deamonDNSPod)

	if (ips == nil) {
		isOk = false
		msg = "Failed to lookup ip in pod for name " + deamonDNSPod
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

func checkMasterApis(dc *models.DeamonClient, urls string) {
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get(toCall)
	return err == nil
}

func checkPodHttpAtoB(dc *models.DeamonClient) {
	// This should fail as we do not have access to this project
	handleCheckStarted(dc)
	var msg string

	isOk := !checkHttp(deamonDNSServiceB + ":8090")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_POD_SERVICE_A_B, IsOk: isOk, Message: msg}
}

func checkPodHttpAtoC(dc *models.DeamonClient) {
	// This should work as we joined this projects
	handleCheckStarted(dc)
	var msg string

	isOk := checkHttp(deamonDNSServiceC + ":8090")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_POD_SERVICE_A_C, IsOk: isOk, Message: msg}
}

func checkHttpHaProxy(dc *models.DeamonClient, publicUrl string) {
	handleCheckStarted(dc)
	var msg string

	isOk := checkHttp(publicUrl + ":80")

	handleCheckFinished(dc, isOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.HTTP_HAPROXY, IsOk: isOk, Message: msg}
}
