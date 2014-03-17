package gotop

import (
	"encoding/json"
	"testing"
)

func TestTotalMemory(t *testing.T) {
	done := make(chan struct{})
	memInfoChan, errc := TotalMemory(done, Delay)
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
	defer func() {
		totalMemoryFile = "/proc/meminfo"
	}()
	done := make(chan struct{})
	memInfoChan, errc := TotalMemory(done, Delay)
	defer func() {
		close(done)
		<-errc
	}()
	for {
		select {
		case <-memInfoChan:
			t.Fatal("Should not return anything")
		case err := <-errc:
			if err == nil {
				t.FailNow()
			}
			return
		}
	}
}
