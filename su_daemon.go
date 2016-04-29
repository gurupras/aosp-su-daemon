package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/alecthomas/kingpin"
)

func write(data, file string) (ret int) {
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

func execv(cmd string, args []string, show_output bool) (ret int) {
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

func Main(args []string) (ret int) {
	switch kingpin.MustParse(app.Parse(args[1:])) {
	case write_cmd.FullCommand():
		debug("write %s %s", write_data, write_file)
		ret = write(*write_data, *write_file)
	case exec_cmd.FullCommand():
		debug("exec %s %s", exec_cmd_bin, exec_cmd_args)
		ret = execv(*exec_cmd_bin, *exec_cmd_args, *output)
	}
	return
}
