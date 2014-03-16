package gotop

import (
	"encoding/json"
	// "fmt"
	"testing"
)

func TestNothing(t *testing.T) {
}

func TestProcessInfo(t *testing.T) {
	done := make(chan struct{})
	processInfoChan, errc := GetProcessInfo(done, Delay)
	for i := 0; ; i = i + 1 {
		if i == 10 {
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
