package main

import (
	"fmt"
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
func TestBase64Gif(t *testing.T) {
	s, err := base64Gif()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
	if len(s) < 1 {
		t.Error("Expected a none empty string")
	}

}

func TestGetLogo(t *testing.T) {
	res, err := getLogo()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}

	if len(res) < 1 {
		t.Error("Expected none empty string")
	}
}

func TestGetType(t *testing.T) {
	res, _ := getLogo()
	ty := getType(res)

	if len(ty) < 1 {
		t.Error("Expected none empty string")
	}
}
