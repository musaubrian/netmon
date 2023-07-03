package main

import (
	"errors"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("COULD NOT LOAD .env\nRENAME `.env.example` to `.env` AND ADJUST THE CONTENTS")
	}

	return nil
}
