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
		t.Errorf("Expected [%s] in template", expected)
	}
}
