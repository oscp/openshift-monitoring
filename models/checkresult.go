package models

type CheckResult struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
	IsOk     bool   `json:"isOk"`
	Message  string `json:"message"`
}
