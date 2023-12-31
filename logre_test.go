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

func TestWriteFatalLog(t *testing.T) {

	e := "Sample error"
	res := WriteFatalLog(e)
	if res != nil {
		t.Errorf("Expected nil got %v", res)
	}
}
