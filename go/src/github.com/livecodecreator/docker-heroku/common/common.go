package common

import "time"

// RaspStatus is
type RaspStatus struct {
	CPU       string    `json:"cpu"`
	Disk      string    `json:"disk"`
	Memory    string    `json:"memory"`
	BootTime  string    `json:"bootTime"`
	Timestamp time.Time `json:"timestamp"`
}

// LastRaspStatus is
var LastRaspStatus RaspStatus
