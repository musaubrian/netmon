package main

import (
	"strings"
	"testing"
)

func TestServiceLocTempl(t *testing.T) {
	location := &ServiceLocation{
		URL: "localhost:8000",
	}

	got, err := serviceLocTempl(location)
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
	str := got.String()
	expected := "Service Location"
	if !strings.Contains(str, expected) {
		t.Errorf("Expected [%s] in template, %s", expected, str)
	}
}

func TestAlertMailTempl(t *testing.T) {
	alert := &Alert{
		LastSpike: Spike{
			T: "12:28",
		},
	}

	res, err := alertMailTempl(alert)
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
	if !strings.Contains(res.String(), alert.LastSpike.T) {
		t.Errorf("Expected [%s] to be in template", alert.LastSpike.T)
	}
}

func TestBackOnlineNotif(t *testing.T) {
	lg := &LastLog{
		Date:     "2023-07-11",
		Time:     "12:28",
		URL:      "localhost:8000",
		Start:    "12:28",
		Duration: "12 seconds",
	}
	str := "An outage occured at"
	res, err := backOnlineNotif(lg)
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}
	if !strings.Contains(res.String(), str) {
		t.Errorf("Expected [%s] to be in template", lg.Date)
	}
}
