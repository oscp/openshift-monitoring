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
	job := getJobStruct(msg)

	h.lastJobId++
	job.JobId = h.lastJobId
	job.JobStatus = models.JOB_RUNNING

	h.jobs[job.JobId] = &job

	// Start job on deamons
	h.jobStart <- job

	// Return ok to UI
	return models.BaseModel{WsType: models.WS_NEW_JOB, Message: job.JobId}
}

func stopJob(h *Hub, msg interface{}) models.BaseModel {
	job := getJobStruct(msg)

	h.jobs[job.JobId].JobStatus = models.JOB_STOPPED
	h.jobStop <- job.JobId

	return models.BaseModel{WsType: models.WS_STOP_JOB}
}

func getJobStruct(msg interface{}) models.Job {
	var job models.Job
	err := mapstructure.Decode(msg, &job)
	if err != nil {
		log.Println("error decoding job", err)
	}
	return job
}