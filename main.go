package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	g "github.com/musaubrian/netmon/gno"
	probing "github.com/prometheus-community/pro-bing"
)

var (
	dRecs   Logr
	records [][]Record
	spikes  []string

	down       bool
	netDownErr string
	alertOnUp  = true

	downTimeStart    time.Time
	savedOutDownTime = false
)

func main() {
	if err := loadEnv(); err != nil {
		WriteFatalLog(err.Error())
		g.Log(g.ERROR, err.Error())
	}
	if err := loadConfig(); err != nil {
		WriteFatalLog(err.Error())
		g.Log(g.ERROR, err.Error())
	}

	today := time.Now()
	todayStr := minimalDate(today.Format(time.RFC850))

	ctx := context.Background()
	timeOutCount := 0

	// Create a ticker that ticks every minute
	ticker := time.NewTicker(5 * time.Minute)

	ngrok_token := Config().NgrokToken

	tunn, err := createNgrokListener(ctx, ngrok_token)
	if err != nil {
		WriteFatalLog(err.Error())
		g.Log(g.ERROR, err.Error())
	}

	// Start up http server
	go Server(ctx, tunn)

	if err := serverLocMail(tunn.URL()); err != nil {
		WriteFatalLog(err.Error())
		g.Log(g.ERROR, err.Error())
	}

	startNetmon(Config().S, timeOutCount, ticker, todayStr, tunn.URL())
}

func startNetmon(s string, tCount int, t *time.Ticker, today string, uri string) {
	// adjust the date every three hours
	dateTicker := time.NewTicker(3 * time.Hour)
	for {
		var r []Record
		maxLat := Config().MaxLat
		down = false
		// initializa new pinger
		pinger, err := probing.NewPinger(s)
		if err != nil {
			err = errors.Join(errors.New("PINGER INITIALIZATION ERR: "), err)
			WriteFatalLog(err.Error())
			g.Log(g.ERROR, err.Error())
		}

		// WINDOWS PRIVILEGES
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		pinger.Timeout = time.Duration(Config().TimeOut * int(time.Millisecond))

		start := time.Now()
		err = pinger.Run()
		if err != nil {
			netDownErr = createLogErr(err.Error(), time.Now())
			down = true
			alertOnUp = false
			if down && !savedOutDownTime {
				downTimeStart = time.Now()
				g.Log(g.WARN, "NETWORK DOWN")
			}
			savedOutDownTime = true
		}

		// ping results
		stats := pinger.Statistics()
		latency := stats.AvgRtt

		if !down && !alertOnUp {
			down = false
			alertOnUp = true
			savedOutDownTime = false

			l := &LastLog{
				URL:      uri,
				Start:    downTimeStart.Format(time.TimeOnly),
				Duration: formatDuration(time.Since(downTimeStart)),
			}

			if err := notifyOnBackOnline(l); err != nil {
				g.Log(g.WARN, err.Error())
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
				g.Log(g.ERROR, err.Error())
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
				g.Log(g.WARN, err.Error())
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
		// Update the Display date every 3hrs
		select {
		case <-t.C:
			records = clearRecords(records)
		case <-dateTicker.C:
			n := time.Now()
			today = minimalDate(n.Format(time.RFC850))
		default:
		}

	}
}

func Config() *NetMonConf {
	return &NetMonConf{
		Email:      os.Getenv("EMAIL"),
		Pwd:        os.Getenv("PASSWORD"),
		NgrokToken: os.Getenv("NGROK_TOKEN"),
		S:          getServerToPing(),
		Port:       getPort(),
		Recipients: getEmails(),
		Department: getDeptName(),
		MaxLat:     getMaxLat(),
		TimeOut:    getPingerTimeout(),
	}
}

func clearRecords(r [][]Record) [][]Record {
	return [][]Record{}
}

func formatDuration(t time.Duration) string {
	v := float64(t) / float64(time.Second)
	if v >= 60.00 {
		v = float64(t) / float64(time.Minute)
		return fmt.Sprintf("%.2f Minutes", v)
	}

	return fmt.Sprintf("%.2f Seconds", v)
}
