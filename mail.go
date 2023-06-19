package main

import (
	"bufio"
	"errors"
	"log"
	"net/smtp"
	"os"
)

func sendMail() error {
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

	message := []byte("To:" + to[0] + "\r\n" + "Subject: Possible downtime detected!\r\n" +
		"\r\n" +
		"You might want to have the other ISP on standby")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpAddr, auth, from, to, message)
	if err != nil {
		return err
	}

	log.Fatal("Email sent successfully!")
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
