package gotop

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var (
	uptimeFile  = "/proc/uptime"
	cpuStatFile = "/proc/stat"
)

type CPUStat struct {
	AverageUtilization float64
	CPUUtilization     []float64
}

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("Could not read file.")
	}
	return string(bytes), nil
}

func numberOfCpus() int {
	numberOfCpus := 0
	cpuStats, err := readFile(cpuStatFile)
	if err != nil {
		return 0
	}
	lines := strings.Split(cpuStats, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Index(lines[i], "cpu") == 0 {
			numberOfCpus = numberOfCpus + 1
		} else {
			break
		}
	}
	// The first line of the file is the average of all CPUs
	return numberOfCpus - 1
}

func getStats(snapshotOne, snapshotTwo string) float64 {
	cpuTimesOne, cpuTimesTwo := strings.Split(snapshotOne[5:], " "), strings.Split(snapshotTwo[5:], " ")
	activeTime, idleTime := 0.0, 0.0
	for i := 0; i < len(cpuTimesOne); i++ {
		time1, _ := strconv.ParseFloat(cpuTimesOne[i], 32)
		time2, _ := strconv.ParseFloat(cpuTimesTwo[i], 32)
		if i != 3 {
			activeTime = activeTime + time2 - time1
		} else {
			// Idle time is the fourth column
			idleTime = +time2 - time1
		}
	}
	activePercentage := 100 * activeTime / (activeTime + idleTime)
	// Return value is truncated to 2 places after decimal
	return float64(int(100*activePercentage)) / 100
}

func cpuSnapshot() ([]string, error) {
	snapshot, err := readFile(cpuStatFile)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(snapshot, "\n"), nil
}

func CPUUsage(cpuStatChan chan CPUStat, interval time.Duration) {
	defer close(cpuStatChan)
	numberOfCpus := numberOfCpus()
	stats := CPUStat{AverageUtilization: 0.0, CPUUtilization: make([]float64, numberOfCpus)}
	var currentSnapshot, previousSnapshot []string
	currentSnapshot, err := cpuSnapshot()
	if err != nil {
		return
	}
	for {
		previousSnapshot = currentSnapshot
		time.Sleep(interval)
		currentSnapshot, err = cpuSnapshot()
		if err != nil {
			return
		}
		stats.AverageUtilization = getStats(currentSnapshot[0], previousSnapshot[0])
		for i := 1; i <= numberOfCpus; i++ {
			stats.CPUUtilization[i-1] = getStats(currentSnapshot[i], previousSnapshot[i])
		}

		cpuStatChan <- stats
	}
}

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
