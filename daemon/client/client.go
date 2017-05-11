package client

import (
	"net"
	"github.com/cenkalti/rpc2"
	"github.com/oscp/openshift-monitoring/models"
	"os"
	"log"
	"time"
	"strings"
	"github.com/oscp/openshift-monitoring/daemon/client/checks"
)

func StartDaemon(h string, dt string, ns string) *rpc2.Client {
	// Local state
	host, _ := os.Hostname()
	d := models.Daemon{Hostname: host,
		Namespace: ns,
		DaemonType: dt,
		StartedChecks: 0,
		FailedChecks:0,
		SuccessfulChecks:0}

	dc := &models.DaemonClient{Daemon: d,
		Quit: make(chan bool),
		ToHub: make(chan models.CheckResult)}

	// Register on hub
	conn, _ := net.Dial("tcp", h)
	dc.Client = rpc2.NewClient(conn)
	dc.Client.Handle("startChecks", func(client *rpc2.Client, checks *models.Checks, reply *string) error {
		startChecks(dc, checks)
		*reply = "ok"
		return nil
	})
	dc.Client.Handle("stopChecks", func(client *rpc2.Client, stop *bool, reply *string) error {
		stopChecks(dc)
		*reply = "ok"
		return nil
	})

	// Start handling from & to hub
	go dc.Client.Run()
	go handleCheckResultToHub(dc)

	registerOnHub(h, dc)

	return dc.Client
}

func StopDaemon(c *rpc2.Client) {
	unregisterOnHub(c)
}

func startChecks(dc *models.DaemonClient, checkConfig *models.Checks) {
	tickExt := time.Tick(time.Duration(checkConfig.CheckInterval) * time.Millisecond)
	tickInt := time.Tick(3 * time.Second)

	log.Println("starting async checks")

	go func() {
		for {
			select {
			case <-dc.Quit:
				HandleChecksStopped(dc)
				return
			case <-tickInt:
				if (checkConfig.MasterApiCheck) {
					go func() {
						HandleCheckStarted(dc)
						isOk, msg := checks.CheckMasterApis(checkConfig.MasterApiUrls)
						HandleCheckFinished(dc, isOk, msg, models.MASTER_API_CHECK)
					}()
				}
				if (checkConfig.EtcdCheck && dc.Daemon.IsMaster()) {
					go func() {
						HandleCheckStarted(dc)
						isOk, msg := checks.CheckEtcdHealth(checkConfig.EtcdIps, checkConfig.EtcdCertPath)
						HandleCheckFinished(dc, isOk, msg, models.ETCD_HEALTH)
					}()
				}
			case <-tickExt:
				if (checkConfig.DnsCheck) {
					go func() {
						HandleCheckStarted(dc)
						isOk, msg := checks.CheckDnsNslookupOnKubernetes()
						HandleCheckFinished(dc, isOk, msg, models.DNS_NSLOOKUP_KUBERNETES)
					}()

					if (dc.Daemon.IsNode() || dc.Daemon.IsMaster()) {
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckDnsServiceNode()
							HandleCheckFinished(dc, isOk, msg, models.DNS_SERVICE_NODE)
						}()
					}

					if (dc.Daemon.IsPod()) {
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckDnsInPod()
							HandleCheckFinished(dc, isOk, msg, models.DNS_SERVICE_POD)
						}()
					}
				}

				if (checkConfig.HttpChecks) {
					if (dc.Daemon.IsPod() && strings.HasSuffix(dc.Daemon.Namespace, "a")) {
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckPodHttpAtoB()
							HandleCheckFinished(dc, isOk, msg, models.HTTP_POD_SERVICE_A_B)
						}()
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckPodHttpAtoC(false)
							HandleCheckFinished(dc, isOk, msg, models.HTTP_POD_SERVICE_A_C)
						}()
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckPodHttpAtoC(true)
							HandleCheckFinished(dc, isOk, msg, models.HTTP_POD_SERVICE_A_C)
						}()
					}

					if (dc.Daemon.IsNode() || dc.Daemon.IsMaster()) {
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckHttpService(false)
							HandleCheckFinished(dc, isOk, msg, models.HTTP_SERVICE_ABC)
						}()
						go func() {
							HandleCheckStarted(dc)
							isOk, msg := checks.CheckHttpService(true)
							HandleCheckFinished(dc, isOk, msg, models.HTTP_SERVICE_ABC)
						}()
					}

					go func() {
						HandleCheckStarted(dc)
						isOk, msg := checks.CheckHttpHaProxy(checkConfig.DaemonPublicUrl, false)
						HandleCheckFinished(dc, isOk, msg, models.HTTP_HAPROXY)
					}()

					go func() {
						HandleCheckStarted(dc)
						isOk, msg := checks.CheckHttpHaProxy(checkConfig.DaemonPublicUrl, true)
						HandleCheckFinished(dc, isOk, msg, models.HTTP_HAPROXY)
					}()
				}
			}
		}
	}()
}

func stopChecks(dc *models.DaemonClient) {
	dc.Quit <- true
}

