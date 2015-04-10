package daemon

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDisk(t *testing.T) {
	done := make(chan struct{})
	// Disk usually needs a much longer delay compared to other metrics
	diskChan, errc := DiskRate(done, 5 * Delay)
	var success bool
	timeout := time.After(10 * Delay)
	defer func() {
		close(done)
		// Necessary to read from error channel to prevent sending goroutine going into deadlock
		<-errc
	}()
	for i := 0; ; i++ {
		if i == 10 {
			return
		}
		select {
		case disk := <-diskChan:
			a, _ := json.Marshal(disk)
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
