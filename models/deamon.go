package models

import "github.com/cenkalti/rpc2"

type Daemon struct {
	Hostname         string
	Namespace        string
	DaemonType       string
	StartedChecks    int
	SuccessfulChecks int
	FailedChecks     int
}

func (d *Daemon) IsMaster() bool {
	return d.DaemonType == "MASTER"
}

func (d *Daemon) IsNode() bool {
	return d.DaemonType == "NODE"
}

func (d *Daemon) IsPod() bool {
	return d.DaemonType == "POD"
}

type DaemonClient struct {
	Daemon Daemon
	Client *rpc2.Client
	Quit   chan bool
	ToHub  chan CheckResult
}
