package gotop

import (
	// "fmt"
	"strconv"
	"strings"
	"time"
)

var (
	memInfoFile = "/proc/meminfo"
	memInfoFields = []string {"MemTotal", "MemFree", "Buffers", "Cached", "SwapTotal", "SwapFree"}
)

type MemInfo struct {
	MemTotal  int
	MemFree   int
	Buffers   int
	Cached    int
	SwapTotal int
	SwapFree  int
}

func MemoryUsage(memInfoChan chan MemInfo, interval time.Duration) {
	defer close(memInfoChan)
	memInfoMap := make(map[string]int)
	for {
		memoryData, err := readFile(memInfoFile)
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
		memInfoChan <- memInfo
		time.Sleep(interval)
	}
}

func getMemInfo(data map[string]int) MemInfo {
	var memInfo MemInfo
	for key, value := range data {
		switch(key) {
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
	indexTwo := strings.LastIndex(line, " ")
	val, _ := strconv.Atoi(line[indexOne:indexTwo])
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