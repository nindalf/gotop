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
		t.FailNow()
	}
	fmt.Println(uptimeDuration)
}

func TestUptimeWrongFile(t *testing.T) {
	uptimeFile = "/proc/wrongfile"
	_, err := Uptime()
	if err == nil {
		t.FailNow()
	}
	uptimeFile  = "/proc/uptime"
}

func TestUpSince(t *testing.T) {
	upSince, err := UpSince()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(upSince)
}

func TestUpSinceWrongFile(t *testing.T) {
	uptimeFile = "/proc/wrongfile"
	_, err := UpSince()
	if err == nil {
		t.FailNow()
	}
	uptimeFile  = "/proc/uptime"
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
	done := make(chan struct{})
	cpuInfoChan, errc := CPUUsage(done, time.Second)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			done <- struct{}{}
		}
		select {
		case cpuInfo := <-cpuInfoChan:
			a, _ := json.Marshal(cpuInfo)
			fmt.Println(string(a))
		case err := <-errc:
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
			return
		}
	}
}

func TestCPUUsageWrongFile(t *testing.T) {
	cpuStatFile = "/proc/wrongfile"
	done := make(chan struct{})
	cpuInfoChan, errc := CPUUsage(done, time.Second)
	for {
		select {
		case cpuInfo := <-cpuInfoChan:
			a, _ := json.Marshal(cpuInfo)
			fmt.Println(string(a))
		case err := <-errc:
			if err == nil {
				t.FailNow()
			}
			return
		}
	}
	cpuStatFile = "/proc/stat"
}
