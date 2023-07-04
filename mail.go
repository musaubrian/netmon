package main

import (
	"log"
	"net/smtp"
	"os"
)

type ServerLocation struct {
	URL string
}

func possibleDowntimeMail() error {
	// Recipient(s) email address(es)
	recipients := getEmails()
	mime := "Content-Type: text/html; charset=utf-8\r\n"

	msg := `Houston we have a problem!

There has been a worrisome development within our domain.
Latencies have soared to unprecedented heights, threatening our network's very existence. 
You must identify the elusive root cause to restore order.

Netmon signing off`

	for _, recipient := range recipients {
		email := []byte("To:" + recipient + "\r\n" + "Subject: Latency Anomaly - Requesting Investigation\r\n" +
			mime +
			msg)

		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}
	log.Println("L1 NOTIFIED")
	return nil
}

/*
Send the server's host IP to concerned parties.

Ideally should only ever happen once when the program is launched
*/
func serverLocMail(uri string) error {
	recipients := getEmails()
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
	from := os.Getenv("email")
	password := os.Getenv("pwd")

	smtpHost := "smtp.gmail.com"
	smtpAddr := smtpHost + ":587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
