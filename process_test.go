package gotop

import (
	"encoding/json"
	// "fmt"
	"testing"
)

func TestNothing(t *testing.T) {
	processIds()
}

func TestProcessData(t *testing.T) {
	done := make(chan struct{})
	processInfoChan, errc := ProcessData(done, testingDelay)
	for i := 0; ; i = i + 1 {
		if i == 3 {
			close(done)
		}
		select {
		case processInfo := <-processInfoChan:
			a, _ := json.Marshal(processInfo)
			t.Log(string(a))
		case err := <-errc:
			if err != nil {
				t.Fatal(err)
			}
			return
		}
	}
}
