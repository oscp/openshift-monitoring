package models

import "github.com/cenkalti/rpc2"

type Daemon struct {
	Hostname         string `json:"hostname"`
	Namespace        string `json:"namespace"`
	DaemonType       string `json:"daemonType"`
	StartedChecks    int    `json:"startedChecks"`
	SuccessfulChecks int    `json:"successfulChecks"`
	FailedChecks     int    `json:"failedChecks"`
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
