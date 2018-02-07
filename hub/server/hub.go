package server

import (
	"github.com/cenkalti/rpc2"
	"github.com/oscp/openshift-monitoring/models"
	"log"
	"net"
	"time"
	"github.com/mitchellh/mapstructure"
)

type Hub struct {
	hubAddr       string
	daemons       map[string]*models.DaemonClient
	currentChecks models.Checks
	result        models.Results
	startChecks   chan models.Checks
	stopChecks    chan bool
	ResetStats    chan bool
	toUi          chan models.BaseModel
	updateStats   bool

	// Temp values between ticks
	successfulSinceTick int
	failedSinceTick     int
}

func NewHub(hubAddr string, masterApiUrls string, daemonPublicUrl string,
	etcdIps string, etcdCertPath string) *Hub {

	return &Hub{
		hubAddr:     hubAddr,
		daemons:     make(map[string]*models.DaemonClient),
		startChecks: make(chan models.Checks),
		stopChecks:  make(chan bool),
		ResetStats:  make(chan bool),
		toUi:        make(chan models.BaseModel, 1000),
		updateStats: false,
		result: models.Results{
			SuccessfulChecksByType: make(map[string]int),
			FailedChecksByType:     make(map[string]int),
			Ticks:                  []models.Tick{},
			Errors:                 []models.Failures{},
		},
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
	statsTicker := time.NewTicker(1 * time.Second)
	toUITicker := time.NewTicker(1 * time.Second)

	// Handle stats
	go func() {
		for {
			select {

			case <-h.ResetStats:
				h.resetStats()
				break

			case <-toUITicker.C:
				// Update checkresults & daemons
				h.toUi <- models.BaseModel{Type: models.CheckResults, Message: h.result}
				h.toUi <- models.BaseModel{Type: models.AllDaemons, Message: h.Daemons()}
				break

			case <-statsTicker.C:
				h.aggregateStats()
				break

			case checks := <-h.startChecks:
				h.updateStats = true
				for _, d := range h.daemons {
					if err := d.Client.Call("startChecks", checks, nil); err != nil {
						log.Println("error starting checks on daemon", err)
					}
				}
				break

			case stop := <-h.stopChecks:
				if stop {
					h.updateStats = false
					for _, d := range h.daemons {
						if err := d.Client.Call("stopChecks", stop, nil); err != nil {
							log.Println("error stopping checks on daemon", err)
						}
					}
				}
				break
			}
		}
	}()

	// Create rpc server for communication with clients
	srv := rpc2.NewServer()
	srv.Handle("register", func(c *rpc2.Client, d *models.Daemon, reply *string) error {
		h.AddDaemon(d, c)
		*reply = "ok"
		return nil
	})
	srv.Handle("unregister", func(cl *rpc2.Client, host *string, reply *string) error {
		h.RemoveDaemon(*host)
		*reply = "ok"
		return nil
	})
	srv.Handle("updateCheckcount", func(cl *rpc2.Client, d *models.Daemon, reply *string) error {
		h.daemons[d.Hostname].Daemon = *d
		*reply = "ok"
		return nil
	})
	srv.Handle("checkResult", func(cl *rpc2.Client, r *models.CheckResult, reply *string) error {
		go h.handleCheckResult(r)
		*reply = "ok"
		return nil
	})
	lis, err := net.Listen("tcp", h.hubAddr)
	srv.Accept(lis)
	if err != nil {
		log.Fatalf("Cannot start rpc2 server: %s", err)
	}
}

func (h *Hub) RemoveDaemon(host string) {
	log.Println("daemon left: ", host)
	delete(h.daemons, host)

	h.toUi <- models.BaseModel{Type: models.DaemonLeft, Message: host}
}

func (h *Hub) AddDaemon(d *models.Daemon, c *rpc2.Client) {
	log.Println("new daemon joined:", d)

	h.daemons[d.Hostname] = &models.DaemonClient{Client: c, Daemon: *d}

	if h.currentChecks.IsRunning {
		// Tell the new daemon to join the checks
		if err := c.Call("startChecks", h.currentChecks, nil); err != nil {
			log.Println("error starting checks on newly joined daemon", err)
		}
	}

	h.toUi <- models.BaseModel{Type: models.NewDaemon, Message: d.Hostname}
}

func (h *Hub) StopChecks() models.BaseModel {
	// Save current state & tell daemons
	h.currentChecks.IsRunning = false
	h.stopChecks <- true

	// Return ok to UI
	return models.BaseModel{Type: models.CurrentChecks, Message: h.currentChecks}
}

func (h *Hub) StartChecks(msg interface{}) models.BaseModel {
	checks := getChecksStruct(msg)

	// Save current state & tell daemons
	checks.IsRunning = true
	h.currentChecks = checks
	h.startChecks <- checks

	// Return ok to UI
	return models.BaseModel{Type: models.CurrentChecks, Message: checks}
}

func getChecksStruct(msg interface{}) models.Checks {
	var checks models.Checks
	err := mapstructure.Decode(msg, &checks)
	if err != nil {
		log.Println("error decoding checks", err)
	}
	return checks
}

func (h *Hub) handleCheckResult(r *models.CheckResult) {
	// Write values from check result to temp values
	if r.IsOk {
		h.result.SuccessfulChecks++
		h.result.SuccessfulChecksByType[r.Type]++
		h.successfulSinceTick++
	} else {
		h.result.FailedChecks++
		h.result.FailedChecksByType[r.Type]++
		h.failedSinceTick++

		h.result.Errors = append(h.result.Errors, models.Failures{
			Date:     time.Now(),
			Type:     r.Type,
			Hostname: r.Hostname,
			Message:  r.Message,
		})
	}
}

func (h *Hub) aggregateStats() {
	// Update global fields
	h.result.StartedChecks = 0
	h.result.FinishedChecks = 0
	for _, d := range h.daemons {
		h.result.StartedChecks += d.Daemon.StartedChecks
		h.result.FinishedChecks += d.Daemon.SuccessfulChecks + d.Daemon.FailedChecks
	}

	if h.failedSinceTick > 0 || h.successfulSinceTick > 0 {
		// Create a new tick out of temp values since last tick
		h.result.Ticks = append(h.result.Ticks, models.Tick{
			FailedChecks:     h.failedSinceTick,
			SuccessfulChecks: h.successfulSinceTick,
		})

		// Prepare for next tick
		h.failedSinceTick = 0
		h.successfulSinceTick = 0
	}
}

func (h *Hub) resetStats() {
	log.Println("resetting stats and current results.")
	h.result = models.Results{
		SuccessfulChecksByType: make(map[string]int),
		FailedChecksByType:     make(map[string]int),
		Ticks:                  []models.Tick{},
		Errors:                 []models.Failures{},
	}

	for _, d := range h.daemons {
		d.Daemon.SuccessfulChecks = 0
		d.Daemon.FailedChecks = 0
		d.Daemon.StartedChecks = 0
	}
}
