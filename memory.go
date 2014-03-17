package gotop

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

var (
	TotalMemoryFile = "/proc/meminfo"
	memInfoFields   = []string{"MemTotal", "MemFree", "Buffers", "Cached", "SwapTotal", "SwapFree"}
)

type MemInfo struct {
	MemTotal  int
	MemFree   int
	Buffers   int
	Cached    int
	SwapTotal int
	SwapFree  int
}

func TotalMemory(done <-chan struct{}, delay time.Duration) (<-chan MemInfo, <-chan error) {
	result := make(chan MemInfo, 1)
	errc := make(chan error)
	var err error
	var memInfo MemInfo
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	go func() {
		defer cleanup()
		for {
			memInfo, err = getMemInfo()
			if err != nil {
				return
			}
			select {
			case result <- memInfo:
			case <-done:
				return
			}
			time.Sleep(delay)
		}
	}()
	return result, errc
}

func getMemInfo() (MemInfo, error) {
	var memInfo MemInfo

	memoryData, err := readFile(TotalMemoryFile)
	if err != nil {
		return memInfo, err
	}
	for _, line := range strings.Split(memoryData, "\n") {
		field := fieldName(line)
		if isMemInfoField(field) {
			value := fieldValue(line)
			switch field {
			case "MemTotal":
				memInfo.MemTotal = value
			case "MemFree":
				memInfo.MemFree = value
			case "Buffers":
				memInfo.Buffers = value
			case "Cached":
				memInfo.Cached = value
			case "SwapTotal":
				memInfo.SwapTotal = value
			case "SwapFree":
				memInfo.SwapFree = value
			}
		}
	}
	return memInfo, nil
}

func fieldName(line string) string {
	index := strings.Index(line, ":")
	if index >= 0 {
		return line[:index]
	}
	return ""
}

func fieldValue(line string) int {
	indexOne := strings.IndexAny(line, "0123456789")
	indexTwo := strings.LastIndexAny(line, "0123456789")
	val, _ := strconv.Atoi(line[indexOne : indexTwo+1])
	return val
}

func isMemInfoField(field string) bool {
	for _, val := range memInfoFields {
		if field == val {
			return true
		}
	}
	return false
}
