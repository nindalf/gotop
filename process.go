package gotop

import (
	// "errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	procDirectory   = "/proc"
	processStatFile = "/proc/%d/stat"
	processStates   = map[string]string{"R": "Running", "S": "Sleeping", "D": "Sleeping-uninterruptable", "Z": "Zombie", "T": "Traced/Stopped"}
)

type ProcessInfo struct {
	Pid   int
	Name  string
	State string
}

func readProcessStatFile(pid int) ([]string, error) {
	path := fmt.Sprintf(processStatFile, pid)
	file, err := readFile(path)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(file, " "), nil
}

func processInfo(pid int) (ProcessInfo, error) {
	pfile, err := readProcessStatFile(pid)
	if err != nil {
		return ProcessInfo{}, err
	}
	name := pfile[1]
	name = name[1 : len(name)-1] // Removes parentheses
	state := processStates[pfile[2]]
	return ProcessInfo{pid, name, state}, nil
}

func processIds() ([]int, error) {
	file, err := os.Open(procDirectory)
	if err != nil {
		return make([]int, 0), err
	}
	fi, _ := file.Readdir(-1)
	pids := make([]int, len(fi))
	var index int
	for _, i := range fi {
		val, err := strconv.Atoi(filepath.Base(i.Name()))
		if err != nil {
			continue
		}
		pids[index] = val
		index = index + 1
	}
	return pids[:index], nil
}

func ProcessData(done <-chan struct{}, interval time.Duration) (<-chan []ProcessInfo, <-chan error) {
	resultChan := make(chan []ProcessInfo)
	errc := make(chan error)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultChan)
	}
	go func() {
		defer cleanup()
		result := make([]ProcessInfo, 1000)
		var index int
		for {
			pids, err := processIds()
			if err != nil {
				return
			}
			index = 0
			for _, pid := range pids {
				pinfo, err := processInfo(pid)
				if err != nil {
					return
				}
				result[index] = pinfo
				index = index + 1
			}
			select {
			case resultChan <- result[:len(pids)]:
			case <-done:
				return
			}
			time.Sleep(interval)
		}
	}()
	return resultChan, errc
}
