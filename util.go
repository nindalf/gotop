package gotop

import (
	"errors"
	"io/ioutil"
)

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("Could not read file.")
	}
	return string(bytes), nil
}