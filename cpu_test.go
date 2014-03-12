package gotop

import (
	"fmt"
	"testing"
)

func TestUptime(t *testing.T) {
	uptimeDuration, err := Uptime()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(uptimeDuration)
}