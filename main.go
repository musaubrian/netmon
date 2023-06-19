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

func minimalDate(d string) string {
	var f string
	u := strings.Split(d, " ")
	f = u[0] + " " + u[1]
	return f
}

func main() {
	today := time.Now()
	todayStr := minimalDate(today.Format(time.RFC850))
	server := "8.8.8.8"

	// Start up http server
	go Server()

	// Create a ticker that ticks every minute
	ticker := time.NewTicker(5 * time.Minute)

	for {
		var r []record
		// initializa new pinger
		pinger, err := probing.NewPinger(server)
		if err != nil {
            log.Fatal("PINGER INITIALIZATION ERR: ",err)
		}
        pinger.Count = 5
		pinger.Timeout = 500 * time.Millisecond

		start := time.Now()
		err = pinger.Run()
		if err != nil {
			if err := sendMail(); err != nil {
				log.Fatal(err)
			}
			log.Println("PINGER ERR: ", err)
		}
		// results from the 5 pings
		stats := pinger.Statistics()
		// fmt.Printf("%+v\n\n", *stats)
		latency := stats.AvgRtt
		if stats.PacketLoss > 25 {
			err := sendMail()
			if err != nil {
				log.Println(err)
			}
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
