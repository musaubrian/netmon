package main

import (
	"fmt"
	"os"
	"time"
)

func WriteFatalLog(e string) error {
	t := time.Now()
	// add a new line to the result
	c := createLogErr(e, t)
	content := c + "\n---\n\n"
	p := "logs" + string(os.PathSeparator) + "fatal"
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o660)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(content)
	return nil
}

func WriteLatenciesLog() error {
	t := time.Now()
	e := fmt.Sprintf("Latencies exceeded %dms\n", getMaxLat())
	content := createLogErr(e, t)
	p := "logs" + string(os.PathSeparator) + "latencies"
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o660)
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
