package gotop

import (
	"encoding/json"
	"fmt"
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
			fmt.Println(string(a))
		case err := <-errc:
			fmt.Println(err)
			return
		}
	}
}
