package main

import (
	"context"
	"os"
	"testing"
)

func TestCreateNGROKListener(t *testing.T) {
	loadEnv()
	c := context.Background()
	token := os.Getenv("ngrok_token")
	_, err := createNgrokListener(c, token)
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
}
