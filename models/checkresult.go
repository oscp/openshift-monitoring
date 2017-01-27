package models

type CheckResult struct {
	Hostname string
	Type     string
	IsOk     bool
	Message  string
}
