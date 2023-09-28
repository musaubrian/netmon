package main

import (
	"log"
	"net/smtp"
)

const MIME = "Content-Type: text/html; charset=utf-8\r\n"

func possibleDowntimeMail(t *Alert) error {
	// Recipient(s) email address(es)
	recipients := Config().Recipients

	body, err := alertMailTempl(t)
	if err != nil {
		return err
	}
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: 🚨 Latency Anomaly - Requesting Investigation\r\n" +
			MIME + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}
	log.Printf("[%s] NOTIFIED ON NETWORK ISSUES\n", Config().Department)
	return nil
}

/*
Send the server's location to concerned parties.

Only ever happen once when the program is launched
*/
func serverLocMail(uri string) error {
	recipients := Config().Recipients
	loc := &ServiceLocation{
		URL: uri,
	}

	body, err := serviceLocTempl(loc)
	if err != nil {
		log.Fatal(err)
	}

	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: Service location\r\n" + MIME + "\r\n" + body.String())
		err := sendMail(recipient, email)
		if err != nil {
			return err
		}
	}

	log.Printf("SHARED NETMON'S LOCATION TO [%s]\n", Config().Department)
	return nil
}

/*
If the network was down and comes back up

Notify the necessary people
*/
func notifyOnBackOnline(lg *LastLog) error {
	recipients := Config().Recipients
	s, err := cleanNetDownErr(netDownErr)
	if err != nil {
		return err
	}

	lg.Date = s[0]
	lg.Time = s[1]

	body, err := backOnlineNotif(lg)
	if err != nil {
		return err
	}
	for _, recipient := range recipients {
		email := []byte("To:" + recipient +
			"\r\nSubject: 📡 Back Online\r\n" +
			MIME + "\r\n" + body.String())
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
