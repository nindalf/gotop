package daemon

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
	resultc := make(chan map[string]ProcessInfo, 1)
	errc := make(chan error)
	pInfoChan := processStats(done, delay)
	pMemChan := processMems(done, delay)
	var err error
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultc)
	}
	go func() {
		defer cleanup()
		if err = checkProcDirectory(); err != nil {
			return
		}
		var pMem map[int]int
		var pInfo map[int]ProcessInfo
		for {
			select {
			case pInfo = <-pInfoChan:
				pMem = <-pMemChan
				result := merge(pInfo, pMem)
				resultc <- result
			case pMem = <-pMemChan:
				pInfo = <-pInfoChan
				result := merge(pInfo, pMem)
				resultc <- result
			case <-done:
				return
			}
		}
	}()
	return resultc, errc
}

func processMems(done <-chan struct{}, delay time.Duration) <-chan map[int]int {
	resultc := make(chan map[int]int, 1)
	go func() {
		defer close(resultc)
		var pids []int
		var result map[int]int
		psize := pagesize()
		for {
			pids = processIds()
			result = make(map[int]int)
			for _, pid := range pids {
				mem := processMem(pid)
				result[pid] = mem
			}
			for key, val := range result {
				result[key] = val * psize
			}
			select {
			case resultc <- result:
			case <-done:
				return
			}
		}
	}()
	return resultc
}

func processStats(done <-chan struct{}, delay time.Duration) <-chan map[int]ProcessInfo {
	resultc := make(chan map[int]ProcessInfo, 1)
	go func() {
		defer close(resultc)
		var pids []int
		var result map[int]ProcessInfo
		prev := make(map[int]pstat)
		var cpu float64
		var timecur, timeprev int64
		for {
			timeprev = timecur
			time.Sleep(delay)
			timecur = time.Now().UnixNano()
			pids = processIds()
			result = make(map[int]ProcessInfo)
			for _, pid := range pids {
				pscur := processStat(pid)
				if psprev, ok := prev[pid]; ok {
					cpu = calcCPU(psprev, pscur, timeprev, timecur)
				} else {
					cpu = 0
				}
				result[pid] = ProcessInfo{pscur.pid, pscur.name, pscur.state, cpu, 0}
				prev[pid] = pscur
			}

			select {
			case resultc <- result:
			case <-done:
				return
			}
		}
	}()
	return resultc
}

func processIds() []int {
	file, err := os.Open(procDirectory)
	if err != nil {
		return make([]int, 0)
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
	return pids[:index]
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

func processStat(pid int) pstat {
	pfile, err := readStatFile(pid)
	if err != nil {
		return pstat{}
	}
	name := pfile[1]
	name = name[1 : len(name)-1] // Removes parentheses
	state := pfile[2]            //processStates[pfile[2]]
	utime, _ := strconv.Atoi(pfile[13])
	stime, _ := strconv.Atoi(pfile[14])
	starttime, _ := strconv.Atoi(pfile[21])
	return pstat{pid, name, state, utime, stime, starttime}
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

func processMem(pid int) int {
	memfile := fmt.Sprintf(processMemFile, pid)
	val, err := readFile(memfile)
	if err != nil {
		return 0
	}
	vals := stringtointslice(val)
	return vals[1]
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
	if len(pInfo) == 0 || len(pMem) == 0 || len(result) == 0 {
		fmt.Println("some length is zero. pinfo, pmem, result", len(pInfo), len(pMem), len(result))
	}
	return result
}

func checkProcDirectory() error {
	memfile := fmt.Sprintf(processMemFile, 1)
	if _, err := readFile(memfile); err != nil {
		return err
	}
	if _, err := readStatFile(1); err != nil {
		return err
	}
	if _, err := os.Open(procDirectory); err != nil {
		return err
	}
	return nil
}
