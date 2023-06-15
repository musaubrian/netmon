package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type record struct {
	Start   time.Time `json:"start"`
	Latency float64   `json:"latency"`
}

type logr struct {
	Day       string     `json:"day"`
	UpdatedAt string     `json:"updated_at"`
	Records   [][]record `json:"records"`
}

func main() {
	var fName string
	var records [][]record
	today := time.Now()
	todayStr := today.Format(time.ANSIC)
	server := "https://8.8.8.8"
	f := strings.Split(todayStr, " ")
	for i := 0; i < 3; i++ {
		fName += "_" + f[i]
	}

	limit := 1 * time.Second
	timeOutErr := "context deadline exceeded"
	client := &http.Client{
		Timeout: limit,
	}

	for {
		var r []record

		start := time.Now()
		_, err := client.Get(server)
		if err != nil {
			if err.Error() == timeOutErr {
				log.Fatal("TIMED OUT:", err)
			} else {
				log.Fatal("ERROR:", err)
			}
		}

		latency := time.Since(start)
		r = append(r,
			record{
				Start: start, Latency: latency.Seconds(),
			})
		records = append(records, r)

		dRecs := logr{
			Day:       todayStr,
			UpdatedAt: time.Now().Format(time.TimeOnly),
			Records:   records,
		}

		go func() {
			err := createLog(fName, dRecs)
			if err != nil {
				log.Println(err)
			}
		}()
		time.Sleep(5 * time.Second)

	}
}

func createLog(fName string, d logr) error {
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
}

func printJSON(d logr) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}

	fmt.Println(string(jsonData))
}
