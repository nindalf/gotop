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
		return 0, errors.New("Could not read file")
	}
	uptime := strings.Split(string(bytes), " ")[0]
	uptimeDuration, err := time.ParseDuration(uptime + "s")
	if err != nil {
		return 0, errors.New("Could not parse uptime")
	}
	return uptimeDuration, nil
}