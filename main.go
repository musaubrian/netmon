package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
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
	limit   = 600 * time.Millisecond
	records [][]record
)

func main() {
	var fName string
	today := time.Now()
	todayStr := today.Format(time.ANSIC)
	server := "https://8.8.8.8"
	timeOutCount := 0
	tRec := []uint16{}

	f := strings.Split(todayStr, " ")
	for i := 0; i < 3; i++ {
		fName += "_" + f[i]
	}
	fName = fName + "_errs"

	client := &http.Client{
		Timeout: limit,
	}
	// Start up http server
	go Server()

	// Create a ticker that ticks every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)

	for {
		var r []record

		start := time.Now()
		_, err := client.Get(server)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				timeOutCount++
			} else {
				log.Println("ERROR:", err)
				timeOutCount++
			}
		}

		latency := time.Since(start)
		r = append(r,
			record{
				Start: start, Latency: uint16(latency.Milliseconds()),
			})
		tRec = append(tRec, uint16(latency.Milliseconds()))
		if timeOutCount == 5 {
			alertOnTimeouts(tRec[len(tRec)-5:])
			// Reset values
			tRec = []uint16{}
			timeOutCount = 0
		}

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
	}
}

func alertOnTimeouts(timeOuts []uint16) {
	fmt.Println(timeOuts)
	averageTimeout := avgLatency(timeOuts)
	fmt.Println(averageTimeout)
	fmt.Printf("Average of 5 timeouts: %d\n", averageTimeout)
	if averageTimeout >= uint16(limit) || averageTimeout <= uint16(limit)+1 {
		fmt.Println("Average of 3 timeouts reached 800!")
	}
}

func avgLatency(l []uint16) uint16 {
	var avg uint16
	var t uint16

	for _, v := range l {
		t += v
	}
	avg = t / uint16(len(l))

	return avg
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
