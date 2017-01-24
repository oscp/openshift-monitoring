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

	go func() {
		for {
			select {
			case <-dc.Quit:
				log.Println("stopped all checks")
				updateChecksCount(dc, true)
				return
			case <-tick:
				if (checks.MasterApiCheck) {
					updateChecksCount(dc, false)
					go checkMasterApis(dc, checks.MasterApiUrl)
				}
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()
}

func stopChecks(dc *models.DeamonClient) {
	dc.Quit <- true
}

func checkMasterApis(dc *models.DeamonClient, urls string) {
 	urlArr := strings.Split(urls, ",")

	oneApiOk := false
	for _,u := range urlArr {
		_, err := http.Get(u)
		if (err == nil) {
			oneApiOk = true
		}
	}

	// Tell the hub about it
	dc.ToHub <- models.CheckResult{Type: models.MASTER_API_CHECK, IsOk: oneApiOk, Message: ""}
}
