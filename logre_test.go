package main

import (
	"strings"
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	e := "Sample error"
	_t := time.Now()
	// Will be in format: 2006-01-02 15:04:05
	formattedT := _t.Format(time.DateTime)
	result := createLogErr(e, _t)

	r := strings.Split(result, " ")
	_fT := r[0] + " " + r[1]
	if formattedT != _fT {
		t.Errorf("Expected [%s] got %s", formattedT, _fT)
	}
}

func TestCreateLogErr(t *testing.T) {
	e := "Sample error"
	_t := time.Now()
	// Will be in format: 2006-01-02 15:04:05
	formattedT := _t.Format(time.DateTime)
	result := createLogErr(e, _t)
	if !strings.Contains(result, e) {
		t.Errorf("Expected [%s] to be in result", e)
	}

	r := strings.Split(result, " ")
	_fT := r[0] + " " + r[1]
	if formattedT != _fT {
		t.Errorf("Expected [%s] got %s", formattedT, _fT)
	}
}

func TestReadNetDownLogs(t *testing.T) {

	lg, err := ReadNetDownLog()
	if err != nil {
		t.Errorf("Expected nil got %v", err)
	}

	if len(lg.Date) < 1 || len(lg.Time) < 1 {
		t.Error("Expected a non empty value")
	}
}
