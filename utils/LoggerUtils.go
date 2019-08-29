package utils

import (
	"fmt"
	"io/ioutil"
)

func ReadLog(file string) (str string, err error) {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str = string(b) // convert content to a 'string'
	return
}
