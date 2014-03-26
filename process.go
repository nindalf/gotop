package gotop

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
#include <unistd.h>
unsigned int get_pagesize(void) {
  return sysconf(_SC_PAGESIZE);
}
*/
import "C"

var (
	procDirectory   = "/proc"
	processStatFile = "/proc/%d/stat"
	processMemFile  = "/proc/%d/statm"
	// processStates   = map[string]string{"R": "Running", "S": "Sleeping", "D": "Sleeping-uninterruptable", "Z": "Zombie", "T": "Traced/Stopped"}
)

type ProcessInfo struct {
	Pid    int
	Name   string
	State  string
	CPU    float64
	Memory int
}

func GetProcessInfo(done <-chan struct{}, delay time.Duration) (<-chan map[string]ProcessInfo, <-chan error) {
	resultChan := make(chan map[string]ProcessInfo, 1)
	errc := make(chan error)
	pInfoChan, errsc := processStats(done, delay)
	pMemChan, errmc := processMems(done, delay)
	var err error
	cleanup := func() {
		fmt.Println(err)
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
				result := merge(pInfo, pMem)
				resultChan <- result
			case pMem = <-pMemChan:
				pInfo = <-pInfoChan
				result := merge(pInfo, pMem)
				resultChan <- result
			case err = <-errsc:
				fmt.Println(err)
				return
			case err = <-errmc:
				fmt.Println(err)
				return
			case <-done:
				return
			}
		}
	}()
	return resultChan, errc
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
		var pids []int
		pids, err = processIds()
		if err != nil {
			return
		}
		var result map[int]int
		var mem int
		psize := pagesize()
		for {
			result = make(map[int]int)
			for _, pid := range pids {
				mem, err = processMem(pid)
				if err != nil {
					return
				}
				result[pid] = mem
			}
			for key, val := range result {
				result[key] = val * psize
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
		var pids []int
		pids, err = processIds()
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
				result[pid] = ProcessInfo{ps.pid, ps.name, ps.state, cpu, 0}
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

type pstat struct {
	pid       int
	name      string
	state     string
	utime     int
	stime     int
	startTime int
}

func processStat(pid int) (pstat, error) {
	pfile, err := readStatFile(pid)
	if err != nil {
		return pstat{}, err
	}
	name := pfile[1]
	name = name[1 : len(name)-1] // Removes parentheses
	state := pfile[2]            //processStates[pfile[2]]
	utime, _ := strconv.Atoi(pfile[13])
	stime, _ := strconv.Atoi(pfile[14])
	starttime, _ := strconv.Atoi(pfile[21])
	return pstat{pid, name, state, utime, stime, starttime}, nil
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

func processMem(pid int) (int, error) {
	memfile := fmt.Sprintf(processMemFile, pid)
	val, err := readFile(memfile)
	if err != nil {
		return 0, err
	}
	vals := stringtointslice(val)
	return vals[1], nil
}

func pagesize() int {
	return int(C.get_pagesize())
}

func merge(pInfo map[int]ProcessInfo, pMem map[int]int) map[string]ProcessInfo {
	result := make(map[string]ProcessInfo)
	for pid, val := range pInfo {
		val.Memory = pMem[pid]
		if val.Memory != 0 {
			result[strconv.Itoa(pid)] = val
		}
	}
	return result
}
