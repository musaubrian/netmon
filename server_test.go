package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCreateNGROKListener(t *testing.T) {
	loadEnv()
	ticker := time.NewTimer(2 * time.Second)
	c := context.Background()
	token := os.Getenv("NGROK_TOKEN")
	_, err := createNgrokListener(c, token)
	<-ticker.C
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
}
