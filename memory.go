package gotop

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

var (
	totalMemoryFile = "/proc/meminfo"
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
	result := make(chan MemInfo)
	errc := make(chan error)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	memInfoMap := make(map[string]int)
	var memoryData string
	go func() {
		defer cleanup()
		for {
			memoryData, err = readFile(totalMemoryFile)
			if err != nil {
				return
			}
			for _, line := range strings.Split(memoryData, "\n") {
				field := fieldName(line)
				if isMemInfoField(field) {
					memInfoMap[field] = fieldValue(line)
				}
			}
			memInfo := getMemInfo(memInfoMap)
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

func getMemInfo(data map[string]int) MemInfo {
	var memInfo MemInfo
	for key, value := range data {
		switch key {
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
	return memInfo
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
