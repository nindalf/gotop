package gotop

import (
	"encoding/json"
	"testing"
)

func TestTotalMemory(t *testing.T) {
	done := make(chan struct{})
	memInfoChan, errc := TotalMemory(done, Delay)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			close(done)
		}
		select {
		case memInfo := <-memInfoChan:
			a, _ := json.Marshal(memInfo)
			t.Log(string(a))
		case err := <-errc:
			if err != nil {
				t.Fatal(err)
			}
			return
		}
	}
}

func TestMemoryUsageWrongFile(t *testing.T) {
	totalMemoryFile = "/proc/wrongfile"
	done := make(chan struct{})
	memInfoChan, errc := TotalMemory(done, Delay)
	for {
		select {
		case memInfo := <-memInfoChan:
			a, _ := json.Marshal(memInfo)
			t.Log(string(a))
		case err := <-errc:
			if err == nil {
				t.FailNow()
			}
			return
		}
	}
	totalMemoryFile = "/proc/meminfo"
}
