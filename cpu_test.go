package gotop

import (
	"encoding/json"
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
	average := sum / float64(len(input))
	return float64(int(average*100)) / 100
}

func TestCPUUsage(t *testing.T) {
	cpuInfoChan := make(chan CPUInfo)
	go CPUUsage(cpuInfoChan, time.Second)
	iterations := 0
	for {
		cpuInfo, ok := <-cpuInfoChan
		iterations = iterations + 1
		if ok == false || iterations > 3 {
			break
		}
		a, _ := json.Marshal(cpuInfo)
		fmt.Println(string(a))
	}
}
