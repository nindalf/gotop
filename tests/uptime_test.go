package gotop

import (
	"testing"
)

func TestUptime(t *testing.T) {
	uptimeDuration, err := Uptime()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(uptimeDuration)
}

func TestUptimeWrongFile(t *testing.T) {
	uptimeFile = "/proc/wrongfile"
	defer func() {
		uptimeFile = "/proc/uptime"
	}()
	_, err := Uptime()
	if err == nil {
		t.FailNow()
	}
}

func TestUpSince(t *testing.T) {
	upSince, err := UpSince()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(upSince)
}

func TestUpSinceWrongFile(t *testing.T) {
	uptimeFile = "/proc/wrongfile"
	defer func() {
		uptimeFile = "/proc/uptime"
	}()
	_, err := UpSince()
	if err == nil {
		t.FailNow()
	}
}
