package main

import (
	"bufio"
	"errors"
	"log"
	"net/smtp"
	"os"
)

func possibleDowntimeMail() error {
	// Recipient email address(es)
	to := []string{"musaubrian45@gmail.com"}

	msg := `Greetings, mortal!

There has been a worrisome development within our domain. The latencies have soared to unprecedented heights, threatening our network's very existence. 
Summon your expertise and investigate this matter with utmost urgency. It could just be a fluke or it could be the start of something that's not good

The fate of our network lays in your hands.

Signed,
netmon`

	for _, recipient := range to {
		email := []byte("To:" + recipient + "\r\n" + "Subject: Latency Anomaly - Requesting Investigation\r\n" +
			"\r\n" +
			msg)

		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}
	log.Println("L1 NOTIFIED")
	return nil
}

func serverLocMail(ip string) error {
	to := []string{"musaubrian45@gmail.com"}

	for _, recipient := range to {

		msg := "Greetings,\r\nI'm up and running at " + ip +
			":8000\r\n" + "I'll notify you if something doesn't seem right\r\n" +
			"\r\nSigned,\r\nnetmon"

		email := []byte("To:" + recipient + "\r\nSubject: Server location\r\n" + "\r\n" + msg)
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}

	log.Println("SERVER LOCATION SHARED")
	return nil
}

func sendMail(to string, msg []byte) error {
	data, err := getDets()
	if err != nil {
		return err
	}
	from := data[0]
	password := data[1]

	smtpHost := "smtp.gmail.com"
	smtpAddr := smtpHost + ":587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpAddr, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

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
