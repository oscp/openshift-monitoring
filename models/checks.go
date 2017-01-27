package models

type Checks struct {
	IsRunning      bool
	CheckInterval  int
	MasterApiCheck bool
	MasterApiUrls  string
	DnsCheck       bool
}
