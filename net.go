package gotop

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	netdir = "/sys/class/net"
	rxfile = "/sys/class/net/%s/statistics/rx_bytes"
	txfile = "/sys/class/net/%s/statistics/tx_bytes"
)

type Net struct {
	Rxrate float64
	Txrate float64
}

func connNames() ([]string, error) {
	file, err := os.Open(netdir)
	if err != nil {
		return make([]string, 0), err
	}
	fi, _ := file.Readdir(-1)
	conns := make([]string, len(fi))
	for index, i := range fi {
		conns[index] = i.Name()
	}
	return conns, nil
}

func numbytes(path string) (float64, error) {
	conns, err := connNames()
	if err != nil {
		return 0, err
	}
	var total int
	for _, conn := range conns {
		p := fmt.Sprintf(path, conn)
		val, err := readFile(p)
		if err != nil {
			return 0, err
		}
		val = strings.Trim(val, " \n")
		// Skip the lo interface. Its the loopback used by local processes
		if val == "lo" {
			continue
		}
		num, _ := strconv.Atoi(val)
		total = total + num
	}
	return float64(total), nil
}

func getrate(prevbytes, curbytes float64, prevtime, curtime int64) float64 {
	// Needs to be converted from nanoseconds to seconds
	timedelta := float64(curtime-prevtime) / 1000000000
	rate := (curbytes - prevbytes) / (timedelta * 1024)
	// Return value is truncated to 2 places after decimal
	return float64(int(100*rate)) / 100
}

func NetRate(done <-chan struct{}, delay time.Duration) (<-chan Net, <-chan error) {
	result := make(chan Net, 1)
	errc := make(chan error)
	var err error
	var net Net
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	go func() {
		defer cleanup()
		var prevrx, currx, prevtx, curtx float64
		var prevtime, curtime int64
		for {
			prevrx = currx
			currx, err = numbytes(rxfile)
			if err != nil {
				return
			}
			prevtx = curtx
			curtx, err = numbytes(txfile)
			if err != nil {
				return
			}
			prevtime = curtime
			curtime = time.Now().UnixNano()
			net.Rxrate = getrate(prevrx, currx, prevtime, curtime)
			net.Txrate = getrate(prevtx, curtx, prevtime, curtime)
			select {
			case result <- net:
			case <-done:
				return
			}
			time.Sleep(delay)
		}
	}()
	return result, errc
}
