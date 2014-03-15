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
	cpuStatFile = "/proc/stat"
)

type CPUInfo struct {
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

func readCPUFile() ([]string, error) {
	snapshot, err := readFile(cpuStatFile)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(snapshot, "\n"), nil
}

func TotalCPU(done <-chan struct{}, interval time.Duration) (<-chan CPUInfo, <-chan error) {
	result := make(chan CPUInfo)
	errc := make(chan error, 1)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	go func() {
		defer cleanup()
		numberOfCpus := numberOfCpus()
		cpuInfo := CPUInfo{AverageUtilization: 0.0, CPUUtilization: make([]float64, numberOfCpus)}
		var currentFile, previousFile []string
		currentFile, err = readCPUFile()
		if err != nil {
			return
		}
		for {
			previousFile = currentFile
			time.Sleep(interval)
			currentFile, err = readCPUFile()
			if err != nil {
				return
			}
			cpuInfo.AverageUtilization = getStats(currentFile[0], previousFile[0])
			for i := 1; i <= numberOfCpus; i++ {
				cpuInfo.CPUUtilization[i-1] = getStats(currentFile[i], previousFile[i])
			}
			select {
			case result <- cpuInfo:
			case <-done:
				return
			}
		}
	}()
	return result, errc
}