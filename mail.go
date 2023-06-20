package main

import (
	"bufio"
	"errors"
	"log"
	"net/smtp"
	"os"
)

func possibleDowntimeMail() error {
	data, err := getDets()
	if err != nil {
		return err
	}
	from := data[0]
	password := data[1]

	// Recipient email address
	to := []string{"musaubrian45@gmail.com"}

	// SMTP server details
	smtpHost := "smtp.gmail.com"
	smtpAddr := smtpHost + ":587"
	msg := `Greetings, mortal!

There has been a worrisome development within our domain. The latencies have soared to unprecedented heights, threatening our network's very existence. 
Summon your expertise and investigate this matter with utmost urgency. It could just be a fluke or it could be the start of something that's not good

The fate of our network lays in your hands.

Signed,
The Superior NetMon`

	email := []byte("To:" + to[0] + "\r\n" + "Subject: Latency Anomaly - Requesting Investigation\r\n" +
		"\r\n" +
		msg)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpAddr, auth, from, to, email)
	if err != nil {
		return err
	}

	log.Println("L1 NOTIFIED")
	return nil
}

func serverLocMail(ip string) error {
	data, err := getDets()
	if err != nil {
		return err
	}
	from := data[0]
	password := data[1]

	// Recipient email address
	to := []string{"musaubrian45@gmail.com"}

	// SMTP server details
	smtpHost := "smtp.gmail.com"
	smtpAddr := smtpHost + ":587"
	msg := "Greetings, mortal!\r\nI am located at http://" + ip + ":8000\r\n\n" + "I'll be on the lookout\r\n" + "\r\n\nSigned,\r\nThe Superior NetMon"
	email := []byte("To:" + to[0] + "\r\n" + "Subject: Server location\r\n" +
		"\r\n" +
		msg)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpAddr, auth, from, to, email)
	if err != nil {
		return err
	}

	log.Println("SERVER LOCATION SHARED")
	return nil
}

func getDets() ([]string, error) {
	var dets []string

	f, err := os.Open(".env")
	if err != nil {
		return dets, errors.Join(errors.New("COULD NOT OPEN FILE.\n"), err)
	}
	sc := bufio.NewScanner(f)

	for sc.Scan() {
		dets = append(dets, sc.Text())
	}

	return dets, nil
}
