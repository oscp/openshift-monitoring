package models

type Checks struct {
	IsRunning       bool   `json:"isRunning"`
	CheckInterval   int    `json:"checkInterval"`
	MasterApiCheck  bool   `json:"masterApiCheck"`
	MasterApiUrls   string `json:"masterApiUrls"`
	DnsCheck        bool   `json:"dnsCheck"`
	HttpChecks      bool   `json:"httpChecks"`
	DaemonPublicUrl string `json:"daemonPublicUrl"`
	EtcdCheck       bool   `json:"etcdCheck"`
	EtcdIps         string `json:"etcdIps"`
	EtcdCertPath    string `json:"etcdCertPath"`
}
