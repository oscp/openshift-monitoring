package models

type Checks struct {
	IsRunning       bool
	CheckInterval   int
	MasterApiCheck  bool
	MasterApiUrls   string
	DnsCheck        bool
	HttpChecks      bool
	DeamonPublicUrl string
	EtcdCheck       bool
	EtcdIps         string
}
