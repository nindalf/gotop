package gotop

import (
	"encoding/json"
	"testing"
)

func TestTotalCPU(t *testing.T) {
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, Delay)
	defer func() {
		close(done)
		// Necessary to read from error channel to prevent sending goroutine going into deadlock
		<-errc
	}()
	for i := 0; ; i = i + 1 {
		if i == 3 {
			return
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
	defer func() {
		totalCPUFile = "/proc/stat"
	}()
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, Delay)
	defer func() {
		close(done)
		<-errc
	}()
	for {
		select {
		case <-cpuInfoChan:
			t.Fatal("Should not return anything")
		case err := <-errc:
			if err == nil {
				t.FailNow()
			}
			return
		}
	}
}
