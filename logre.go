package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WriteNetworkDownErrs(e string, t time.Time) {
	c := createLogErr(e, t)
	content := c + "\n"

	f, err := os.OpenFile("./logs/network_down", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o660)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(content)
}

func WriteFatalErrs(e string) {
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
