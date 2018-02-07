package models

import "time"

type Results struct {
	SuccessfulChecks int `json:"successfulChecks"`
	FailedChecks     int `json:"failedChecks"`
	StartedChecks    int `json:"startedChecks"`
	FinishedChecks   int `json:"finishedChecks"`

	SuccessfulChecksByType map[string]int `json:"successfulChecksByType"`
	FailedChecksByType     map[string]int `json:"failedChecksByType"`

	Ticks  []Tick     `json:"ticks"`
	Errors []Failures `json:"failures"`
}

type Tick struct {
	SuccessfulChecks int `json:"successfulChecks"`
	FailedChecks     int `json:"failedChecks"`
}

type Failures struct {
	Date     time.Time `json:"date"`
	Hostname string    `json:"hostname"`
	Type     string    `json:"type"`
	Message  string    `json:"message"`
}

type CheckResult struct {
	Hostname string
	Type     string
	IsOk     bool
	Message  string
}
