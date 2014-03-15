package gotop

import (
	"errors"
	"strings"
	"time"
)

var (
	uptimeFile  = "/proc/uptime"
)

func Uptime() (time.Duration, error) {
	uptimeString, err := readFile(uptimeFile)
	if err != nil {
		return 0, err
	}
	uptime := strings.Split(uptimeString, " ")[0]
	uptimeDuration, err := time.ParseDuration(uptime + "s")
	if err != nil {
		return 0, errors.New("Could not parse uptime.")
	}
	return uptimeDuration, nil
}

func UpSince() (time.Time, error) {
	duration, err := Uptime()
	if err != nil {
		return time.Unix(0, 0).UTC(), errors.New("Could not get uptime.")
	}
	return time.Now().Add(-1 * duration), nil
}
