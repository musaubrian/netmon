package main

import (
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

func getLogo() (string, error) {
	var logo string
	loc := "./web/static/"
	defLogo := loc + "netmon.png"

	c, err := os.ReadDir(loc)
	if err != nil {
		return logo, err
	}

	for _, v := range c {
		if strings.Contains(v.Name(), "logo") {
			logo = loc + v.Name()
		}
	}
	if len(logo) < 1 {
		logo = defLogo

	}
	return logo, err
}

func getType(file string) string {
	t := strings.Split(file, ".")
	return t[2]
}
