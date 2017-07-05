package server

import (
	"github.com/cenkalti/rpc2"
	"github.com/oscp/openshift-monitoring/models"
	"log"
	"net"
	"time"
)

type Hub struct {
	hubAddr       string
	daemons       map[string]*models.DaemonClient
	currentChecks models.Checks
	checkResults  []models.CheckResult
	startChecks   chan models.Checks
	stopChecks    chan bool
	toUi          chan models.BaseModel
}

func NewHub(hubAddr string, masterApiUrls string, daemonPublicUrl string,
	etcdIps string, etcdCertPath string) *Hub {

	return &Hub{
		hubAddr:      hubAddr,
		daemons:      make(map[string]*models.DaemonClient),
		startChecks:  make(chan models.Checks),
		stopChecks:   make(chan bool),
		toUi:         make(chan models.BaseModel, 1000),
		checkResults: []models.CheckResult{},
		currentChecks: models.Checks{
			CheckInterval:   5000,
			MasterApiUrls:   masterApiUrls,
			DaemonPublicUrl: daemonPublicUrl,
			MasterApiCheck:  true,
			HttpChecks:      true,
			DnsCheck:        true,
			EtcdCheck:       true,
			EtcdIps:         etcdIps,
			EtcdCertPath:    etcdCertPath,
			IsRunning:       false},
	}
}

func (h *Hub) Daemons() []models.Daemon {
	r := []models.Daemon{}
	for _, d := range h.daemons {
		r = append(r, d.Daemon)
	}
	return r
}

func (h *Hub) Serve() {
	go handleChecksStart(h)
	go handleChecksStop(h)
	go updateUI(h)

	srv := rpc2.NewServer()
	srv.Handle("register", func(c *rpc2.Client, d *models.Daemon, reply *string) error {
		// Save client for talking to him later
		daemonJoin(h, d, c)
		*reply = "ok"
		return nil
	})
	srv.Handle("unregister", func(cl *rpc2.Client, host *string, reply *string) error {
		daemonLeave(h, *host)

		*reply = "ok"
		return nil
	})
	srv.Handle("updateCheckcount", func(cl *rpc2.Client, d *models.Daemon, reply *string) error {
		h.daemons[d.Hostname].Daemon = *d
		*reply = "ok"
		return nil
	})
	srv.Handle("checkResult", func(cl *rpc2.Client, r *models.CheckResult, reply *string) error {
		h.checkResults = append(h.checkResults, *r)
		*reply = "ok"
		return nil
	})
	lis, err := net.Listen("tcp", h.hubAddr)
	srv.Accept(lis)
	if err != nil {
		log.Fatalf("Cannot start rpc2 server: %s", err)
	}
}

func updateUI(h *Hub) {
	// UI cannot handle each checkresult individually, so we aggrgate and send them each second
	tick := time.Tick(1 * time.Second)
	go func() {
		for {
			select {
			case <-tick:
				// Update checkresults
				h.toUi <- models.BaseModel{Type: models.CHECK_RESULTS, Message: h.checkResults}
				h.checkResults = []models.CheckResult{}

				// Update deamons
				h.toUi <- models.BaseModel{Type: models.ALL_DAEMONS, Message: h.Daemons()}
			}
		}
	}()
}

func handleChecksStart(h *Hub) {
	for {
		var checks models.Checks = <-h.startChecks
		for _, d := range h.daemons {
			if err := d.Client.Call("startChecks", checks, nil); err != nil {
				log.Println("error starting checks on daemon", err)
			}
		}
	}
}

func handleChecksStop(h *Hub) {
	for {
		var stop bool = <-h.stopChecks

		if stop {
			for _, d := range h.daemons {
				if err := d.Client.Call("stopChecks", stop, nil); err != nil {
					log.Println("error stopping checks on daemon", err)
				}
			}
		}
	}
}
