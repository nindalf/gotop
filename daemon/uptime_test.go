package daemon

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

func TestUpSince(t *testing.T) {
	upSince, err := UpSince()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(upSince)
}
