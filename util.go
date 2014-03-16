package gotop

import (
	"errors"
	"fmt"
	"io/ioutil"
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
