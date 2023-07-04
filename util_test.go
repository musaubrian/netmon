package main

import "testing"

func TestLoadEnv(t *testing.T) {
	err := loadEnv()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
}
