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
	_, err := Uptime()
	if err == nil {
		t.FailNow()
	}
	uptimeFile  = "/proc/uptime"
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
	_, err := UpSince()
	if err == nil {
		t.FailNow()
	}
	uptimeFile  = "/proc/uptime"
}
