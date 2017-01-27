package models

import "github.com/cenkalti/rpc2"

type Deamon struct {
	Hostname         string
	DeamonType       string
	StartedChecks    int
	SuccessfulChecks int
	FailedChecks     int
}

func (d *Deamon) IsMaster() bool {
	return d.DeamonType == "MASTER"
}

func (d *Deamon) IsNode() bool {
	return d.DeamonType == "NODE"
}

func (d *Deamon) IsPod() bool {
	return d.DeamonType == "POD"
}

type DeamonClient struct {
	Deamon Deamon
	Client *rpc2.Client
	Quit   chan bool
	ToHub  chan CheckResult
}
