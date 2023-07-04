package main

import "testing"

func TestLoadConfig(t *testing.T) {
	err := loadConfig()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
}
