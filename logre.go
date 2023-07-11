package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// last log written in the network_down log file
type LastLog struct {
	Date string `json:"date"`
	Time string `json:"time"`
	URL  string
}

var mu sync.Mutex

func WriteNetworkDownLog(e string, t time.Time) {
	mu.Lock()
	defer mu.Unlock()

	c := createLogErr(e, t)
	content := c + "\n"

	// Overwrite the contents everytime, making the file only ever one line long
	f, err := os.OpenFile("./logs/network_down", os.O_WRONLY|os.O_CREATE, 0o660)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(content)
}

func WriteFatalLog(e string) {
	t := time.Now()
	// add a new line to the result
	c := createLogErr(e, t)
	content := c + "\n---\n\n"
	f, err := os.OpenFile("./logs/fatal", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o660)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(content)
}

func WriteLatenciesLog() error {
	t := time.Now()
	e := fmt.Sprintf("Latencies exceeded %dms\n", getMaxLat())
	content := createLogErr(e, t)
	f, err := os.OpenFile("./logs/latencies", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

func createLogErr(err string, t time.Time) string {
	var logErr string

	cleanTime := formatTime(t)
	logErr = cleanTime + " " + err
	return logErr
}

func formatTime(t time.Time) string {
	formattedTime := t.Format(time.DateTime)
	return formattedTime
}

func ReadNetDownLog() (*LastLog, error) {
	var l string

	f, err := os.Open("./logs/network_down")
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		l = sc.Text()
	}
	arrL := strings.Split(l, " ")
	lg := &LastLog{
		Date: arrL[0],
		Time: arrL[1],
	}

	return lg, err
}
