package gotop

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	//Delay between samples
	Delay = 500 * time.Millisecond
)

type Systeminfo struct {
	Sysname  string
	Nodename string
	Release  string
	Version  string
	Machine  string
	CPUModel string
	NumCPU   int
	Memory   string
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
	memory := float64(memInfo.MemTotal) / (1024 * 1024)
	memstr := strconv.FormatFloat(memory, 'f', 2, 64) + "GB"
	return Systeminfo{sysname, nodename, release, version, machine, model, numCPU, memstr}
}

func (s Systeminfo) String() string {
	return fmt.Sprintf("Sysname: \t%s\nNodename: \t%s\nRelease: \t%s\nVersion: \t%s\nMachine: \t%s\nModel: \t\t%s\nCPU cores: \t%d\nTotal memory:\t%s",
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

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errorString := fmt.Sprintf("Could not read file %s .", filename)
		return "", errors.New(errorString)
	}
	return string(bytes), nil
}

func stringtointslice(input string) []int {
	// Get rid of extra spaces
	re, _ := regexp.Compile(" +")
	input = re.ReplaceAllLiteralString(input, " ")
	input = strings.Trim(input, " \n")

	temp := strings.Split(input, " ")
	output := make([]int, len(temp))
	var index int
	for _, val := range temp {
		valint, err := strconv.Atoi(val)
		if err == nil {
			output[index] = valint
			index = index + 1
		}
	}
	return output[:index]
}

func getrate(prevval, curval float64, prevtime, curtime int64) float64 {
	// Needs to be converted from nanoseconds to seconds
	timedelta := float64(curtime-prevtime) / 1000000000
	rate := (curval - prevval) / (timedelta * 1024)
	// Return value is truncated to 2 places after decimal
	return float64(int(100*rate)) / 100
}
