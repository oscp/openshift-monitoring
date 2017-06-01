package checks

import (
	"net"
	"log"
	"net/http"
	"crypto/tls"
	"strings"
	"os/exec"
	"regexp"
	"strconv"
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
	var msg string
	out, err := exec.Command("bash", "-c", "ntpstat").Output()
	if err != nil {
		msg = "Could not check ntpd status: " + err.Error()
		log.Println(msg)
		return false, msg
	}

	isOk := strings.Contains(string(out), "time correct")
	if (!isOk) {
		msg = "Time is not correct on the server or ntpd not running"
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
	log.Println("Checking access to:", toCall)
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

func isLvsSizeOk(stdOut string, okSize int) bool {
	// Examples
	// 42.10  8.86   docker-pool
	// 13.63  8.93   lv_fast_registry_pool
	num := regexp.MustCompile("(\\d+\\.\\d+)")

	checksOk := 0
	for _, nr := range num.FindAllString(stdOut, -1) {
		i, err := strconv.ParseFloat(nr, 64)
		if (err != nil) {
			log.Print("Unable to parse int:", nr)
			return false
		}

		if (i < float64(okSize)) {
			checksOk++
		} else {
			log.Println("LVM pool size exceeded okSize:", i)
		}
	}

	return checksOk == 2
}