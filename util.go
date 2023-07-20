package main

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("COULD NOT LOAD .env\nRENAME `.env.example` to `.env` AND ADJUST THE CONTENTS")
	}

	return nil
}

// Format date to: `day, date-month-year`
func minimalDate(d string) string {
	var f string
	u := strings.Split(d, " ")
	f = u[0] + " " + u[1]
	return f
}

func cleanNetDownErr(n string) ([]string, error) {
	var s []string
	if len(n) > 1 {
		s = strings.Split(n, " ")
	} else {
		return s, errors.New("Empty value, cannot be split")
	}
	return s, nil
}

// Contents of the gif as a base64 string
func base64Gif() (string, error) {
	var b64 string
	g, err := os.ReadFile("./web/static/calltronix.gif")
	if err != nil {
		return b64, err
	}
	b64 = base64.StdEncoding.EncodeToString(g)
	return b64, err
}
