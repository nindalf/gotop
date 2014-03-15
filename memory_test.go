package gotop

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMemoryUsage(t *testing.T) {
	done := make(chan struct{})
	memInfoChan, errc := MemoryUsage(done, time.Second)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			done <- struct{}{}
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
	memInfoFile = "/proc/wrongfile"
	done := make(chan struct{})
	memInfoChan, errc := MemoryUsage(done, time.Second)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			done <- struct{}{}
		}
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
	memInfoFile   = "/proc/meminfo"
}
