package gotop

import (
	"fmt"
	"testing"
	"time"
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

func average(input []float64) float64 {
	var sum float64
	for i := 0; i < len(input); i++ {
		sum = sum + input[i]
	}
	average := sum/float64(len(input))
	return float64(int(average * 100)) / 100
}

func TestCPUUsage(t *testing.T) {
	cpuStatChan := make(chan CPUStat)
	go CPUUsage(cpuStatChan, 2*time.Second)
	for cpuStats := range cpuStatChan {
		fmt.Println(cpuStats.AverageUtilization, cpuStats.CPUUtilization, average(cpuStats.CPUUtilization))
	}
}
