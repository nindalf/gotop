package gotop

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestMemoryUsage(t *testing.T) {
	memInfoChan := make(chan MemInfo)
	go MemoryUsage(memInfoChan, time.Second)
	iterations := 0
	for {
		memInfo, ok := <-memInfoChan
		if ok == false || iterations > 3 {
			break
		}
		a, _ := json.Marshal(memInfo)
		fmt.Println(string(a))
	}
}
