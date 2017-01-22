package server

import (
	"github.com/SchweizerischeBundesbahnen/openshift-monitoring/models"
	"log"
	"github.com/cenkalti/rpc2"
	"github.com/mitchellh/mapstructure"
)

func deamonLeave(h *Hub, host string) {
	log.Println("deamon left: ", host)
	delete(h.deamons, host)

	h.toUi <- models.BaseModel{WsType: models.WS_DEAMON_LEFT, Message: host}
}

func deamonJoin(h *Hub, d *models.Deamon, c *rpc2.Client) {
	log.Println("new deamon joined:", d)

	h.deamons[d.Hostname] = models.DeamonClient{Client:c, Deamon: *d}

	h.toUi <- models.BaseModel{WsType: models.WS_NEW_DEAMON, Message: d.Hostname}
}

func newJob(h *Hub, msg interface{}) models.BaseModel {
	var job models.Job
	err := mapstructure.Decode(msg, &job)
	if err != nil {
		panic(err)
	}

	h.lastJobId++
	job.JobId = h.lastJobId
	h.jobs = append(h.jobs, job)

	h.toDeamons <- job

	// Return ok to UI
	return models.BaseModel{WsType: models.WS_NEW_JOB, Message: job.JobId}
}