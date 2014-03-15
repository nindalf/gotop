package gotop

import (
	"encoding/json"
	"testing"
)

func TestTotalCPU(t *testing.T) {
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, testingDelay)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			close(done)
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
	cpuInfoChan, errc := TotalCPU(done, testingDelay)
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
