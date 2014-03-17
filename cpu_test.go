package gotop

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTotalCPU(t *testing.T) {
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, Delay)
	var success bool
	timeout := time.After(2*Delay)
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
			success = true
		case err := <-errc:
			if err != nil {
				t.Fatal(err)
			}
			return
		case <-timeout:
			if success == false {
				t.Fatal("No result. Goroutine hanging.")
			}
		}
	}
}

func TestTotalCPUWrongFile(t *testing.T) {
	totalCPUFile = "/proc/wrongfile"
	done := make(chan struct{})
	cpuInfoChan, errc := TotalCPU(done, Delay)
	var success bool
	timeout := time.After(2*Delay)
	defer func() {
		totalCPUFile = "/proc/stat"
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
		case <-timeout:
			if success == false {
				t.Fatal("No result. Goroutine hanging.")
			}
		}
	}
}
