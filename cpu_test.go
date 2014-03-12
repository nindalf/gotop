package gotop

import (
	"fmt"
	"testing"
)

func TestUptime(t *testing.T) {
	uptimeDuration, err := Uptime()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(uptimeDuration)
}

func TestUpSince(t *testing.T) {
	upSince, err := UpSince()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(upSince)
}