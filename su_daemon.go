package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

//export Execv
func Execv(cmd string, args []string, show_output bool) (ret int) {
	if output, err := exec.Command(cmd, args...).Output(); err != nil {
		if show_output {
			fmt.Fprintln(os.Stderr, err)
		}
		ret = -1
	} else {
		if show_output {
			fmt.Printf("%s", output)
		}
	}
	return
}

//export Execv1
func Execv1(cmd string, args string, show_output bool) (ret int) {
	return Execv(cmd, SliceIt(args), show_output)
}
