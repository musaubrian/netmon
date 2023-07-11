package main

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type NetMonConf struct {
	// For authentication
	Email      string
	Pwd        string
	NgrokToken string

	S          string // Server to ping
	Port       int
	Recipients []string
	MaxLat     int
	TimeOut    int // Maximum pinger timeout
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

var (
	dRecs     Logr
	records   [][]Record
	spikes    []string
	down      bool
	alertOnUp = true
)

func main() {
	if err := loadEnv(); err != nil {
		WriteFatalLog(err.Error())
		log.Fatal(err)
	}
	if err := loadConfig(); err != nil {
		WriteFatalLog(err.Error())
		log.Fatal(err)
	}
	ctx := context.Background()
	today := time.Now()
	todayStr := minimalDate(today.Format(time.RFC850))
	timeOutCount := 0

	// Create a ticker that ticks every minute
	ticker := time.NewTicker(5 * time.Minute)

	ngrok_token := Config().NgrokToken

	tunn, err := createNgrokListener(ctx, ngrok_token)
	if err != nil {
		WriteFatalLog(err.Error())
		log.Fatal(err)
	}
	// Start up http server
	go Server(ctx, tunn)

	if err := serverLocMail(tunn.URL()); err != nil {
		WriteFatalLog(err.Error())
		log.Fatal(err)
	}

	startNetmon(Config().S, timeOutCount, todayStr, ticker, tunn.URL())
}

func startNetmon(s string, tCount int, today string, t *time.Ticker, uri string) {
	for {
		var r []Record
		maxLat := Config().MaxLat
		down = false
		// initializa new pinger
		pinger, err := probing.NewPinger(s)
		if err != nil {
			err = errors.Join(errors.New("PINGER INITIALIZATION ERR: "), err)
			WriteFatalLog(err.Error())
			log.Fatal(err)
		}

		// WINDOWS PRIVILEGES
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		pinger.Timeout = time.Duration(Config().TimeOut * int(time.Millisecond))

		start := time.Now()
		err = pinger.Run()
		if err != nil {
			// Write to file instead of cluttering stdout
			WriteNetworkDownLog(err.Error(), time.Now())
			down = true
			alertOnUp = false
		}

		// ping results
		stats := pinger.Statistics()
		latency := stats.AvgRtt

		if !down && !alertOnUp {
			down = false
			alertOnUp = true
			if err := notifyOnBackOnline(uri); err != nil {
				log.Println(err)
			}
		}
		if int(latency.Milliseconds()) >= maxLat {
			spikes = append(spikes, start.Format(time.TimeOnly))
			tCount++
		}

		// only send alert if more than 3 timeouts have occurred
		if tCount >= 3 {
			tCount = 0
			if err := WriteLatenciesLog(); err != nil {
				log.Fatal(err)
			}
			alert := &Alert{
				MaxLat: maxLat,
				LastSpike: Spike{
					T:   start.Format(time.TimeOnly),
					Lat: uint16(latency.Milliseconds()),
				},
			}
			err := possibleDowntimeMail(alert)
			if err != nil {
				log.Println(err)
			}
			// Sleep for 10 minutes after alerting
			// You have 10 minutes max to find the solution to the issue
			// before it continues
			time.Sleep(10 * time.Minute)
		}

		r = append(r,
			Record{
				Start:   start,
				Latency: uint16(latency.Milliseconds()),
			})

		records = append(records, r)

		dRecs = Logr{
			Day:       today,
			UpdatedAt: time.Now().Format(time.TimeOnly),
			Records:   records,
		}

		// clear records every 5 minutes
		select {
		case <-t.C:
			records = clearRecords(records)
		default:
		}

		// Send two pings per second
		time.Sleep(500 * time.Millisecond)
	}
}

func Config() *NetMonConf {
	return &NetMonConf{
		Email:      os.Getenv("email"),
		Pwd:        os.Getenv("pwd"),
		NgrokToken: os.Getenv("ngrok_token"),
		S:          getServerToPing(),
		Port:       getPort(),
		Recipients: getEmails(),
		MaxLat:     getMaxLat(),
		TimeOut:    getPingerTimeout(),
	}
}

func clearRecords(r [][]Record) [][]Record {
	return [][]Record{}
}
