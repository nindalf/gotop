package gotop

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTotalCPU(t *testing.T) {
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, time.Second)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			done <- struct{}{}
		}
		select {
		case cpuInfo := <-cpuInfoChan:
			a, _ := json.Marshal(cpuInfo)
			t.Log(string(a))
		case err := <-errc:
			if err != nil {
				t.Fatal(err)
			}
			return
		}
	}
}

func TestTotalCPUWrongFile(t *testing.T) {
	totalCPUFile = "/proc/wrongfile"
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, time.Second)
	for {
		select {
		case cpuInfo := <-cpuInfoChan:
			a, _ := json.Marshal(cpuInfo)
			t.Log(string(a))
		case err := <-errc:
			if err == nil {
				t.FailNow()
			}
			return
		}
	}
	totalCPUFile = "/proc/stat"
}
