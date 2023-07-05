package main

import (
	"os"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	expected := &NetMonConf{
		Email:      os.Getenv("email"),
		Pwd:        os.Getenv("pwd"),
		S:          getServerToPing(),
		Port:       getPort(),
		Recipients: getEmails(),
		MaxLat:     getMaxLat(),
		AlertMsg:   getAlertMsg(),
	}

	res := Config()

	if expected.Email != res.Email {
		t.Errorf("Expected %s got %s", expected.Email, res.Email)
	}
	if expected.Pwd != res.Pwd {
		t.Errorf("Expected %s got %s", expected.Pwd, res.Pwd)
	}
	if expected.S != res.S {
		t.Errorf("Expected %s got %s", expected.S, res.S)
	}
	if expected.Port != res.Port {
		t.Errorf("Expected %d got %d", expected.Port, res.Port)
	}
	if expected.MaxLat != res.MaxLat {
		t.Errorf("Expected %d got %d", expected.MaxLat, res.MaxLat)
	}
	if expected.AlertMsg != res.AlertMsg {
		t.Errorf("Expected %s got %s", expected.AlertMsg, res.AlertMsg)
	}
}

func TestClearRecords(t *testing.T) {
	var rc []record
	var rcs [][]record

	rc = append(rc, record{
		Start:   time.Now(),
		Latency: 100,
	})
	for i := 0; i < 3; i++ {
		rcs = append(rcs, rc)
	}
	if len(rcs) != 3 {
		t.Errorf("Expected 3 values got %d", len(rcs))
	}
	rcs = clearRecords(rcs)
	if len(rcs) > 1 {
		t.Errorf("Expected 0 got %d", len(rcs))
	}
}
