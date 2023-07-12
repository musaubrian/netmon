package main

import (
	"errors"
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
