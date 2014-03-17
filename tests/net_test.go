package gotop

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestConnections(t *testing.T) {
	conns, err := connNames()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Conns: ", conns)
}

func TestBytes(t *testing.T) {
	rx, err := numbytes(rxfile)
	if err != nil || rx == 0 {
		t.Fatal(err, rx)
	}
	t.Log("rx bytes: ", rx)
	tx, err := numbytes(txfile)
	if err != nil || tx == 0 {
		t.Fatal(err, tx)
	}
	t.Log("tx bytes: ", tx)
}

func TestNet(t *testing.T) {
	done := make(chan struct{})
	netChan, errc := NetRate(done, 10*Delay)
	var success bool
	timeout := time.After(2 * Delay)
	defer func() {
		close(done)
		// Necessary to read from error channel to prevent sending goroutine going into deadlock
		<-errc
	}()
	for i := 0; ; i = i + 1 {
		if i == 30 {
			return
		}
		select {
		case net := <-netChan:
			a, _ := json.Marshal(net)
			t.Log(string(a))
			fmt.Println(string(a))
			success = true
		case err := <-errc:
			if err != nil {
				t.Fatal(err)
			}
			return
		case <-timeout:
			if success == false {
				t.Fatal("No result. Goroutine hanging.")
			}
		}
	}
}
