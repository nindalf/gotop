package gotop

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"syscall"
	"time"
)

var (
	//Delay between samples
	Delay = 500 * time.Millisecond
)

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errorString := fmt.Sprintf("Could not read file %s .", filename)
		return "", errors.New(errorString)
	}
	return string(bytes), nil
}

type Systeminfo struct {
	Sysname  string
	Nodename string
	Release  string
	Version  string
	Machine  string
	CPUModel string
	NumCPU   int
	Memory string
}

func (s Systeminfo) String() string {
	return fmt.Sprintf("Sysname: %s\nNodename: %s\nRelease: %s\nVersion: %s\nMachine: %s\nModel: %s\nnumCPU: %d\nTotal memory:%s",
		s.Sysname, s.Nodename, s.Release, s.Version, s.Machine, s.CPUModel, s.NumCPU, s.Memory)	
}

func charToStr(input [65]int8) string {
	out := make([]byte, len(input))
	for i, val := range input {
		if val == 0 {
			break
		}
		out[i] = byte(val)
	}
	return string(out)
}	

func Sysinfo() Systeminfo {
	var uts syscall.Utsname
	syscall.Uname(&uts)
	sysname := charToStr(uts.Sysname)
	nodename := charToStr(uts.Nodename)
	release := charToStr(uts.Release)
	version := charToStr(uts.Version)
	machine := charToStr(uts.Machine)
	model := cpuModel()
	numCPU := numberOfCpus()
	memInfo, _ := getMemInfo()
	memory := float64(memInfo.MemTotal)/(1024*1024)
	memstr := strconv.FormatFloat(memory, 'f', 2, 64) + "GB"
	return Systeminfo{sysname, nodename, release, version, machine, model, numCPU, memstr}
}

