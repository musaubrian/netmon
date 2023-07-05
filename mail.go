package main

import (
	"log"
	"net/smtp"
)

type ServerLocation struct {
	URL string
}
type Spike struct {
	T   string
	Lat uint16
}

type Alert struct {
	MaxLat    int
	Message   string
	LastSpike Spike
}

func possibleDowntimeMail(t *Alert) error {
	// Recipient(s) email address(es)
	recipients := Config().Recipients
	t.Message = Config().AlertMsg

	mime := "Content-Type: text/html; charset=utf-8\r\n"
	body, err := alertMailTempl(t)
	if err != nil {
		return err
	}
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: 🚨 Latency Anomaly - Requesting Investigation\r\n" +
			mime + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}
	log.Println("L1 NOTIFIED")
	return nil
}

/*
Send the server's location to concerned parties.

Ideally should only ever happen once when the program is launched
*/
func serverLocMail(uri string) error {
	recipients := Config().Recipients
	loc := &ServerLocation{
		URL: uri,
	}

	body, err := serverLocTempl(loc)
	if err != nil {
		log.Fatal(err)
	}

	mime := "Content-Type: text/html; charset=utf-8\r\n"
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: Server location\r\n" + mime + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}

	log.Println("SERVER LOCATION SHARED")
	return nil
}

func sendMail(to string, msg []byte) error {
	from := Config().Email
	password := Config().Pwd

	smtpHost := "smtp.gmail.com"
	smtpAddr := smtpHost + ":587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
