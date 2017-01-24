package models

import "github.com/cenkalti/rpc2"

type Deamon struct {
	Hostname    string
	DeamonType  string
	ChecksCount int
}

type DeamonClient struct {
	Deamon Deamon
	Client *rpc2.Client
}
