package daemon

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProcessInfo(t *testing.T) {
	done := make(chan struct{})
	processInfoChan, errc := GetProcessInfo(done, Delay)
	var success bool
	timeout := time.After(2 * Delay)
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
		case processInfo := <-processInfoChan:
			a, _ := json.Marshal(processInfo["1"])
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
