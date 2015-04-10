package daemon

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

func NetRate(done <-chan struct{}, delay time.Duration) (<-chan Net, <-chan error) {
	resultc := make(chan Net, 1)
	errc := make(chan error)
	var err error
	var net Net
	cleanup := func() {
		errc <- err
		close(errc)
		close(resultc)
	}
	go func() {
		defer cleanup()
		var prevrx, currx, prevtx, curtx float64
		var prevtime, curtime int64
		if err = checkNetFiles(); err != nil {
			return
		}
		for {
			prevrx = currx
			currx, _ = numbytes(rxfile)
			prevtx = curtx
			curtx, _ = numbytes(txfile)
			prevtime = curtime
			curtime = time.Now().UnixNano()
			net.Rxrate = getrate(prevrx, currx, prevtime, curtime)
			net.Txrate = getrate(prevtx, curtx, prevtime, curtime)
			select {
			case resultc <- net:
			case <-done:
				return
			}
			time.Sleep(delay)
		}
	}()
	return resultc, errc
}

func checkNetFiles() error {
	if _, err := numbytes(rxfile); err != nil {
		return err
	}
	if _, err := numbytes(txfile); err != nil {
		return err
	}
	return nil
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
