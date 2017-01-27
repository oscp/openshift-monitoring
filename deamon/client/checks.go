package client

import (
	"net/http"
	"log"
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"time"
	"strings"
)

func startChecks(dc *models.DeamonClient, checks *models.Checks) {
	tick := time.Tick(1 * time.Second)

	log.Println("starting checks")

	go func() {
		for {
			select {
			case <-dc.Quit:
				log.Println("stopped checks")
				return
			case <-tick:
				if (checks.MasterApiCheck) {
					go checkMasterApis(dc, checks.MasterApiUrls)
				}
			}
		}
	}()
}

func stopChecks(dc *models.DeamonClient) {
	dc.Quit <- true
}

func checkMasterApis(dc *models.DeamonClient, urls string) {
	handleCheckStarted(dc)
 	urlArr := strings.Split(urls, ",")

	oneApiOk := false
	var msg string
	for _,u := range urlArr {
		_, err := http.Get(u)
		if (err == nil) {
			oneApiOk = true
		} else {
			msg += u + " is not reachable. ";
		}
	}

	handleCheckFinished(dc, oneApiOk)

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.MASTER_API_CHECK, IsOk: oneApiOk, Message: msg}
}
