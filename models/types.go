package models

// Websocket types
const (
	NewDaemon  = "NEW_DAEMON"
	AllDaemons = "ALL_DAEMONS"
	DaemonLeft = "DAEMON_LEFT"

	CurrentChecks = "CURRENT_CHECKS"
	StartChecks   = "START_CHECKS"
	StopChecks    = "STOP_CHECKS"
	CheckResults  = "CHECK_RESULTS"
)

// Check types
const (
	MasterApiCheck        = "MASTER_API_CHECK"
	DnsNslookupKubernetes = "DNS_NSLOOKUP_KUBERNETES"
	DnsServiceNode        = "DNS_SERVICE_NODE"
	DnsServicePod         = "DNS_SERVICE_POD"
	HttpPodServiceAB      = "HTTP_POD_SERVICE_A_B"
	HttpPodServiceAC      = "HTTP_POD_SERVICE_A_C"
	HttpHaProxy           = "HTTP_HAPROXY"
	HttpServiceABC        = "HTTP_SERVICE_ABC"
	EtcdHealth            = "ETCD_HEALTH"
)
