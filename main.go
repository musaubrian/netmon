package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type record struct {
	Start   time.Time `json:"start"`
	Latency uint16    `json:"latency"`
}

type logr struct {
	Day       string     `json:"day"`
	UpdatedAt string     `json:"updated_at"`
	Records   [][]record `json:"records"`
}

var (
	dRecs   logr
	records [][]record
)

func main() {
	if err := loadEnv(); err != nil {
		log.Fatal(err)
	}
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	today := time.Now()
	todayStr := minimalDate(today.Format(time.RFC850))
	server := getServerToPing()
	timeOutCount := 0

	// Create a ticker that ticks every minute
	ticker := time.NewTicker(5 * time.Minute)

	ngrok_token := os.Getenv("ngrok_token")

	tunn, err := createNgrokListener(ctx, ngrok_token)
	if err != nil {
		log.Fatal(err)
	}
	// Start up http server
	go Server(ctx, tunn)

	if err := serverLocMail(tunn.URL()); err != nil {
		log.Fatal(err)
	}

	for {
		var r []record
		// initializa new pinger
		pinger, err := probing.NewPinger(server)
		if err != nil {
			log.Fatal("PINGER INITIALIZATION ERR: ", err)
		}

		// WINDOWS PRIVILEGES
		if runtime.GOOS == "windows" {
			pinger.SetPrivileged(true)
		}

		pinger.Timeout = 500 * time.Millisecond

		start := time.Now()
		err = pinger.Run()
		if err != nil {
			log.Println("PINGER ERR: ", err)
		}

		// ping results
		stats := pinger.Statistics()
		latency := stats.AvgRtt

		if int(latency.Milliseconds()) >= 500 {
			timeOutCount++
		}

		// only send alert if more than 3 timeouts have occurred
		if timeOutCount >= 3 {
			timeOutCount = 0
			err := possibleDowntimeMail()
			if err != nil {
				log.Println(err)
			}
			// Sleep for 10 minutes after alerting
			// You have 10 minutes max to find the solution to the issue
			// before it continues
			time.Sleep(10 * time.Minute)
		}

		r = append(r,
			record{
				Start:   start,
				Latency: uint16(latency.Milliseconds()),
			})

		records = append(records, r)

		dRecs = logr{
			Day:       todayStr,
			UpdatedAt: time.Now().Format(time.TimeOnly),
			Records:   records,
		}

		// clear records every 5 minutes
		select {
		case <-ticker.C:
			clearRecords()
		default:
		}

		// time.Sleep(500 * time.Millisecond)
	}
}

// Format date to: `day, date-month-year`
func minimalDate(d string) string {
	var f string
	u := strings.Split(d, " ")
	f = u[0] + " " + u[1]
	return f
}

func clearRecords() {
	records = [][]record{}
}

/*func createLog(fName string, d logr) error {
	f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0o660)
	if err != nil {
		return errors.Join(errors.New("COULD NOT CREATE FILE:"), err)
	}
	defer f.Close()
	jsonData, err := json.Marshal(d)
	if err != nil {
		return errors.Join(errors.New("MARSHALLING FAILED:"), err)
	}
	f.WriteString(string(jsonData))
	fmt.Println("File written to")
	return nil
}*/
