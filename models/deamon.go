package models

import "github.com/cenkalti/rpc2"

type Deamon struct {
	Hostname         string
	DeamonType       string
	StartedChecks    int
	SuccessfulChecks int
	FailedChecks     int
}

type DeamonClient struct {
	Deamon Deamon
	Client *rpc2.Client
	Quit   chan bool
	ToHub  chan CheckResult
}
