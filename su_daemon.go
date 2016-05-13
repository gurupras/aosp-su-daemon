package main

import (
	"bufio"
	"fmt"
	"os"
)

//export Write
func Write(data, file string) (ret int) {
	var f_raw *os.File
	var err error

	if f_raw, err = os.OpenFile(file, os.O_WRONLY, 0); err != nil {
		fmt.Fprintln(os.Stderr, err)
		ret = -1
		return
	}
	defer f_raw.Close()

	writer := bufio.NewWriter(f_raw)
	defer writer.Flush()

	if ret, err = writer.Write([]byte(data)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		ret = -1
		return
	}
	return
}
