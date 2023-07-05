package main

import (
	"strings"
	"testing"
)

func TestServerLocTempl(t *testing.T) {
	location := &ServerLocation{
		URL: "localhost:8000",
	}

	got, err := serverLocTempl(location)
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
