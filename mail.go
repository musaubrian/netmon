package main

import (
	"log"
	"net/smtp"
)

type Static struct {
	Data string // base64 encoding of the gif
}

type ServiceLocation struct {
	URL string
	G   Static
}

type Spike struct {
	T   string
	Lat uint16
}

type Alert struct {
	MaxLat    int
	LastSpike Spike
	G         Static
}

func possibleDowntimeMail(t *Alert) error {
	// Recipient(s) email address(es)
	recipients := Config().Recipients

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
	log.Println("L1 NOTIFIED ON NETWORK ISSUES")
	return nil
}

/*
Send the server's location to concerned parties.

Ideally should only ever happen once when the program is launched
*/
func serverLocMail(uri string, g *Static) error {
	recipients := Config().Recipients
	loc := &ServiceLocation{
		URL: uri,
		G:   *g,
	}

	body, err := serviceLocTempl(loc)
	if err != nil {
		log.Fatal(err)
	}

	mime := "Content-Type: text/html; charset=utf-8\r\n"
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: Service location\r\n" + mime + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}

	log.Println("SHARED NETMON'S LOCATION")
	return nil
}

/*
If the network was down and comes back up

Notify the necessary people
*/
func notifyOnBackOnline(uri string, g *Static) error {
	recipients := Config().Recipients
	s, err := cleanNetDownErr(netDownErr)
	if err != nil {
		return err
	}
	lg := &LastLog{
		Date: s[0],
		Time: s[1],
		URL:  uri,
		G:    *g,
	}

	mime := "Content-Type: text/html; charset=utf-8\r\n"
	body, err := backOnlineNotif(lg)
	if err != nil {
		return err
	}
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: 📡 Back Online\r\n" +
			mime + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}
	log.Println("BACK ONLINE")
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
