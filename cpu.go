package gotop

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

var (
	totalCPUFile = "/proc/stat"
	cpuInfoFile  = "/proc/cpuinfo"
)

// Average utilization is the average of the elements of CPUUtilization.
// Each element of CPUUtilization corresponds to a CPU core.
type CPUInfo struct {
	AverageUtilization float64
	CPUUtilization     []float64
}

//
func TotalCPU(done <-chan struct{}, delay time.Duration) (<-chan CPUInfo, <-chan error) {
	result := make(chan CPUInfo, 1)
	errc := make(chan error)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	go func() {
		defer cleanup()
		var cur, prev []string
		for {
			prev = cur
			cur, err = readCPUFile()
			if err != nil {
				return
			}
			cpuInfo := getCPUInfo(prev, cur)
			select {
			case result <- cpuInfo:
			case <-done:
				return
			}
			time.Sleep(delay)
		}
	}()
	return result, errc
}

func getCPUInfo(prev, cur []string) CPUInfo {
	numberOfCpus := numberOfCpus()
	cpuInfo := CPUInfo{AverageUtilization: 0.0, CPUUtilization: make([]float64, numberOfCpus)}

	// This is the first time the query is happening
	if len(prev) == 0 {
		return cpuInfo
	}

	cpuInfo.AverageUtilization = getStats(cur[0], prev[0])
	for i := 1; i <= numberOfCpus; i++ {
		cpuInfo.CPUUtilization[i-1] = getStats(cur[i], prev[i])
	}
	return cpuInfo
}

func getStats(current, previous string) float64 {
	// Should replace hardcoded value 5. If cpu levels are absurdly high, its because of this
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

func readCPUFile() ([]string, error) {
	snapshot, err := readFile(totalCPUFile)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(snapshot, "\n"), nil
}

func cpuModel() string {
	info, _ := readFile(cpuInfoFile)
	for _, line := range strings.Split(info, "\n") {
		splitline := strings.Split(line, ":")
		field := strings.Trim(splitline[0], " \t")
		if strings.EqualFold(field, "model name") {
			return strings.Trim(splitline[1], " \t")
		}
	}
	return ""
}
