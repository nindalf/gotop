package gotop

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

var (
	totalCPUFile = "/proc/stat"
)

type CPUInfo struct {
	AverageUtilization float64
	CPUUtilization     []float64
}

func TotalCPU(done <-chan struct{}, delay time.Duration) (<-chan CPUInfo, <-chan error) {
	result := make(chan CPUInfo)
	errc := make(chan error)
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
		var cur, prev []string
		cur, err = readCPUFile()
		if err != nil {
			return
		}
		for {
			prev = cur
			time.Sleep(delay)
			cur, err = readCPUFile()
			if err != nil {
				return
			}
			cpuInfo.AverageUtilization = getStats(cur[0], prev[0])
			for i := 1; i <= numberOfCpus; i++ {
				cpuInfo.CPUUtilization[i-1] = getStats(cur[i], prev[i])
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

func numberOfCpus() int {
	numberOfCpus := 0
	cpuStats, err := readCPUFile()
	if err != nil {
		return 0
	}
	for i := 0; i < len(cpuStats); i++ {
		if strings.Index(cpuStats[i], "cpu") == 0 {
			numberOfCpus = numberOfCpus + 1
		} else {
			break
		}
	}
	// The first line of the file is the average of all CPUs
	return numberOfCpus - 1
}

func getStats(current, previous string) float64 {
	// start :=
	prev, cur := strings.Split(previous[5:], " "), strings.Split(current[5:], " ")
	activeTime, idleTime := 0.0, 0.0
	for i := 0; i < len(cur); i++ {
		time1, _ := strconv.ParseFloat(prev[i], 32)
		time2, _ := strconv.ParseFloat(cur[i], 32)
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
	snapshot, err := readFile(totalCPUFile)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(snapshot, "\n"), nil
}
