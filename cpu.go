package gotop

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"strings"
	"time"
)

var (
	uptimeFile = "/proc/uptime"
)


func Uptime() (time.Duration, error){
	bytes, err := ioutil.ReadFile(uptimeFile)
	if err != nil {
		return 0, errors.New("Could not read uptime file.")
	}
	uptime := strings.Split(string(bytes), " ")[0]
	uptimeDuration, err := time.ParseDuration(uptime + "s")
	if err != nil {
		return 0, errors.New("Could not parse uptime.")
	}
	return uptimeDuration, nil
}

func UpSince() (time.Time, error) {
	duration, err := Uptime()
	if err != nil {
		return time.Unix(0, 0), errors.New("Could not get uptime.")
	}
	return time.Now().Add(-1 * duration), nil
}