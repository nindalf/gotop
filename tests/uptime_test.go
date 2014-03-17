package gotop

import (
	"github.com/nindalf/gotop"
	"testing"
)

func TestUptime(t *testing.T) {
	uptimeDuration, err := gotop.Uptime()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(uptimeDuration)
}

func TestUpSince(t *testing.T) {
	upSince, err := gotop.UpSince()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(upSince)
}
