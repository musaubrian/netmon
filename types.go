package main

import "time"

type NetMonConf struct {
	// For authentication
	Email      string
	Pwd        string
	NgrokToken string

	S          string // Server to ping
	Port       int
	Recipients []string
	MaxLat     int

	// Maximum pinger timeout
	// How long to wait for a response before ignoring that ping's results
	// If a ping exceeds this, its result defaults to 0
	TimeOut int
}

type Record struct {
	Start   time.Time `json:"start"`
	Latency uint16    `json:"latency"`
}

type Logr struct {
	Day       string     `json:"day"`
	UpdatedAt string     `json:"updated_at"`
	Records   [][]Record `json:"records"`
}

// The last downtime error
type LastLog struct {
	Date     string `json:"date"`
	Time     string `json:"time"`
	URL      string
	Start    string // Time net went down
	Duration string // How long net was down for
}

type ServiceLocation struct {
	URL string
}

type Spike struct {
	T   string
	Lat uint16
}

type Alert struct {
	MaxLat    int
	LastSpike Spike
}
