package gotop

import (
	// "errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	procDirectory   = "/proc"
	processStatFile = "/proc/%d/stat"
	processMemFile  = "/proc/%d/status"
	processStates   = map[string]string{"R": "Running", "S": "Sleeping", "D": "Sleeping-uninterruptable", "Z": "Zombie", "T": "Traced/Stopped"}
)

type ProcessInfo struct {
	Pid    int
	Name   string
	State  string
	CPU    float64
	Memory int
}

type pstat struct {
	pid       int
	name      string
	state     string
	utime     int
	stime     int
	startTime int
	rss       int
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
		val, err := strconv.Atoi(i.Name())
		if err != nil {
			continue
		}
		pids[index] = val
		index = index + 1
	}
	return pids[:index], nil
}

func readStatFile(pid int) ([]string, error) {
	path := fmt.Sprintf(processStatFile, pid)
	file, err := readFile(path)
	if err != nil {
		return make([]string, 0), err
	}
	return strings.Split(file, " "), nil
}

func processStat(pid int) (pstat, error) {
	pfile, err := readStatFile(pid)
	if err != nil {
		return pstat{}, err
	}
	name := pfile[1]
	name = name[1 : len(name)-1] // Removes parentheses
	state := processStates[pfile[2]]
	utime, _ := strconv.Atoi(pfile[13])
	stime, _ := strconv.Atoi(pfile[14])
	starttime, _ := strconv.Atoi(pfile[21])
	rss, _ := strconv.Atoi(pfile[23])
	return pstat{pid, name, state, utime, stime, starttime, rss}, nil
}

func calcCPU(prevps, curps pstat, prevtime, curtime int64) float64 {
	// Time is in nanoseconds. Needs to be converted to jiffies.
	timedelta := float64(curtime-prevtime) / 10000000
	usercpu := float64(curps.utime-prevps.utime) / timedelta
	systemcpu := float64(curps.stime-prevps.stime) / timedelta
	usage := (usercpu + systemcpu) * 100
	// Return value is truncated to 2 places after decimal
	return float64(int(usage*100)) / 100
}

func processStats(done <-chan struct{}, delay time.Duration) (<-chan map[int]ProcessInfo, <-chan error) {
	resultChan := make(chan map[int]ProcessInfo, 1)
	errc := make(chan error, 1)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultChan)
	}
	go func() {
		defer cleanup()
		pids, err := processIds()
		if err != nil {
			return
		}
		var result map[int]ProcessInfo
		prev := make(map[int]pstat)
		var ps pstat
		var cpu float64
		var curtime, prevtime int64
		for {
			prevtime = curtime
			time.Sleep(delay)
			curtime = time.Now().UnixNano()
			if err != nil {
				return
			}
			result = make(map[int]ProcessInfo)
			for _, pid := range pids {
				ps, err = processStat(pid)
				if err != nil {
					return
				}
				if psprev, ok := prev[pid]; ok {
					cpu = calcCPU(psprev, ps, prevtime, curtime)
				} else {
					cpu = 0
				}
				result[pid] = ProcessInfo{ps.pid, ps.name, ps.state, cpu, ps.rss / (256)}
				prev[pid] = ps
			}

			select {
			case resultChan <- result:
			case <-done:
				return
			}
		}
	}()
	return resultChan, errc
}

func processMem(pid int) (int, error) {
	return 0, nil
}

func processMems(done <-chan struct{}, delay time.Duration) (<-chan map[int]int, <-chan error) {
	resultChan := make(chan map[int]int, 1)
	errc := make(chan error)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultChan)
	}
	go func() {
		defer cleanup()
		pids, err := processIds()
		if err != nil {
			return
		}
		var result map[int]int
		var mem int
		for {
			result = make(map[int]int)
			for _, pid := range pids {
				mem, err = processMem(pid)
				if err != nil {
					return
				}
				result[pid] = mem
			}
			select {
			case resultChan <- result:
			case <-done:
				return
			}
		}
	}()
	return resultChan, errc
}

func merge(pInfo map[int]ProcessInfo, pMem map[int]int) map[int]ProcessInfo {
	for pid, val := range pInfo {
		val.Memory = pMem[pid]
		pInfo[pid] = val
	}
	return pInfo
}

func GetProcessInfo(done <-chan struct{}, delay time.Duration) (<-chan map[int]ProcessInfo, <-chan error) {
	resultChan := make(chan map[int]ProcessInfo, 1)
	errc := make(chan error)
	pInfoChan, errsc := processStats(done, delay)
	pMemChan, errmc := processMems(done, delay)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultChan)
	}
	go func() {
		defer cleanup()
		var pMem map[int]int
		var pInfo map[int]ProcessInfo
		for {
			select {
			case pInfo = <-pInfoChan:
				pMem = <-pMemChan
				pInfo = merge(pInfo, pMem)
			case pMem = <-pMemChan:
				pInfo = <-pInfoChan
				pInfo = merge(pInfo, pMem)
			case err = <-errsc:
				return
			case err = <-errmc:
				return
			}
			select {
			case resultChan <- pInfo:
			case <-done:
				return
			}

		}
	}()
	return resultChan, errc
}
