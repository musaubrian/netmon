package main

import (
	"log"
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
	today := time.Now()
	todayStr := minimalDate(today.Format(time.RFC850))
	server := "8.8.8.8"
	timeOutCount := 0

	// Start up http server
	go Server()

	// Create a ticker that ticks every minute
	ticker := time.NewTicker(5 * time.Minute)

	ip := getIP()
	err := serverLocMail(ip)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var r []record
		// initializa new pinger
		pinger, err := probing.NewPinger(server)
		if err != nil {
			log.Fatal("PINGER INITIALIZATION ERR: ", err)
		}

		// WINDOWS PRIVILEDGES
		// pinger.SetPrivileged(true)

		pinger.Timeout = 500 * time.Millisecond

		start := time.Now()
		err = pinger.Run()
		if err != nil {
			log.Println("PINGER ERR: ", err)
		}

		// ping results
		stats := pinger.Statistics()
		latency := stats.AvgRtt
		if stats.PacketLoss > 40 || int(latency.Milliseconds()) >= 300 {
			timeOutCount++
		}

		// only send alert if more than 5 timeouts have occurred
		if timeOutCount >= 3 {
			err := possibleDowntimeMail()
			if err != nil {
				log.Println(err)
			}
			timeOutCount = 0
		}

		r = append(r,
			record{
				Start: start, Latency: uint16(latency.Milliseconds()),
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
		time.Sleep(1 * time.Second)

	}
}

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
