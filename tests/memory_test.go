package gotop

import (
	"encoding/json"
	"github.com/nindalf/gotop"
	"testing"
	"time"
)

func TestTotalMemory(t *testing.T) {
	done := make(chan struct{})
	memInfoChan, errc := gotop.TotalMemory(done, gotop.Delay)
	var success bool
	timeout := time.After(2 * gotop.Delay)
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