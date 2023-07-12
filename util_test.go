package main

import (
	"testing"
	"time"
)

func TestLoadEnv(t *testing.T) {
	err := loadEnv()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
}
func TestMinimalDate(t *testing.T) {
	unformatted := "Monday, 02-Jan-06 15:04:05 MST"
	expected := "Monday, 02-Jan-06"
	got := minimalDate(unformatted)
	if expected != got {
		t.Errorf("Expected [%s] got [%s]", expected, got)
	}
}

func TestCleanNetDownErr(t *testing.T) {
	e := createLogErr("Some error", time.Now())
	_, err := cleanNetDownErr(e)
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}

}
