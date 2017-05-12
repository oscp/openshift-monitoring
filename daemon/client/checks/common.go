package checks

import (
	"net"
	"log"
	"net/http"
	"crypto/tls"
	"strings"
	"strconv"
	"os/exec"
)

const (
	daemonDNSEndpoint = "daemon.ose-mon-a.endpoints.cluster.local"
	daemonDNSServiceA = "daemon.ose-mon-a.svc.cluster.local"
	daemonDNSServiceB = "daemon.ose-mon-b.svc.cluster.local"
	daemonDNSServiceC = "daemon.ose-mon-c.svc.cluster.local"
	daemonDNSPod = "daemon"
	kubernetesIP = "172.30.0.1"
)


func CheckExternalSystem(url string) (bool, string) {
	isOk := checkHttp(url)

	var msg string
	if (!isOk) {
		msg = "Call to " + url + " failed"
	}

	return isOk, msg
}

func CheckNtpd() (bool, string) {
	isOk := false
	var msg string
	out, err := exec.Command("bash", "-c", "ntpstat").Output()
	if err != nil {
		msg = "Could not check ntpd status: " + err.Error()
		log.Println(msg)
		return isOk, msg
	}

	isOk = strings.Contains(string(out), "time correct")
	if (!isOk) {
		msg = "Time is not correct on the server or ntpd not running"
	}
	return isOk, msg
}

func CheckOpenFileCount() (bool, string) {
	isOk := false
	var msg string
	out, err := exec.Command("bash", "-c", "cat /proc/sys/fs/file-nr | cut -f1").Output()
	if err != nil {
		msg = "Could not evaluate open file count: " + err.Error()
		log.Println(msg)
		return isOk, msg
	}

	nr, err := strconv.Atoi(string(out))

	if (err != nil) {
		msg = "Could not parse output to integer: " + string(out)
		return isOk, msg
	}

	if (nr < 200000) {
		isOk = true
	}

	if (!isOk) {
		msg = "Open files are higher than 200'000 files!"
	}
	return isOk, msg
}

func getIpsForName(n string) []net.IP {
	ips, err := net.LookupIP(n)
	if (err != nil) {
		log.Println("failed to lookup ip for name ", n)
		return nil
	}
	return ips
}

func checkHttp(toCall string) bool {
	if (strings.HasPrefix(toCall, "https")) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(toCall)
		if (err != nil) {
			log.Println("error in http check: ", err)
		} else {
			resp.Body.Close()
		}
		return err == nil
	} else {
		resp, err := http.Get(toCall)
		if (err != nil) {
			log.Println("error in http check: ", err)
		} else {
			resp.Body.Close()
		}
		return err == nil
	}
}

func getEndpoint(slow bool) string {
	if (slow) {
		return "slow"
	} else {
		return "fast"
	}
}