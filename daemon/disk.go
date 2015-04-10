package daemon

import (
	"strconv"
	"strings"
	"time"
)

var (
	diskstatsFile = "/proc/diskstats"
	blocksizeFile = "/sys/block/sda/queue/physical_block_size"
)

type Disk struct {
	Read  float64
	Write float64
}

func DiskRate(done <-chan struct{}, delay time.Duration) (<-chan Disk, <-chan error) {
	result := make(chan Disk, 1)
	errc := make(chan error)
	var err error
	var disk Disk
	cleanup := func() {
		errc <- err
		close(errc)
		close(result)
	}
	go func() {
		defer cleanup()
		var prevwrite, curwrite, prevread, curread float64
		var prevtime, curtime int64
		bsize := float64(blocksize())
		for {
			prevwrite, prevread = curwrite, curread
			curwrite, curread, err = readDiskStats()
			if err != nil {
				return
			}
			prevtime = curtime
			curtime = time.Now().UnixNano()
			disk.Read = getrate(bsize*prevread, bsize*curread, prevtime, curtime)
			disk.Write = getrate(bsize*prevwrite, bsize*curwrite, prevtime, curtime)
			select {
			case result <- disk:
			case <-done:
				return
			}
			time.Sleep(delay)
		}
	}()
	return result, errc
}

func blocksize() int {
	val, err := readFile(blocksizeFile)
	if err != nil {
		// Default block size
		return 512
	}
	valint, _ := strconv.Atoi(strings.Trim(val, " \n"))
	return valint
}

func readDiskStats() (float64, float64, error) {
	val, err := readFile(diskstatsFile)
	if err != nil {
		return 0, 0, err
	}
	sdaline := strings.Index(val, "sda ")
	val = strings.Split(val[sdaline:], "\n")[0]
	valint := stringtointslice(val)
	sectorsRead := float64(valint[5])
	sectorsWritten := float64(valint[9])
	return sectorsWritten, sectorsRead, nil
}
