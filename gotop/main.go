package main

import (
	"fmt"

	"github.com/nindalf/gotop/daemon"
)

func main() {
	fmt.Println(daemon.UpSince())
}
