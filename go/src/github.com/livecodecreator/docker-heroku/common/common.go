package common

import (
	"os"
	"time"
)

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

// EnvStruct is
type EnvStruct struct {
	SlackToken   string
	SlackChannel string
}

// Env is
var Env = EnvStruct{
	SlackToken:   os.Getenv("SLACK_TOKEN"),
	SlackChannel: os.Getenv("SLACK_CHANNEL"),
}
