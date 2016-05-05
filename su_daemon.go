package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/shlex"
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
func Execv(cmd string, args []string, shell bool) (ret int, stdout string, stderr string) {
	var buf_stdout, buf_stderr bytes.Buffer
	var err error
	var command *exec.Cmd

	if shell == true {
		args = append([]string{cmd}, args...)
		argstring := "-c '" + strings.Join(args, " ") + "'"
		args, err = shlex.Split(argstring)
		cmd = ShellPath
	}

	// Create a string to log for the command that we're running
	var cmd_string string = cmd + " "
	for i, arg := range args {
		cmd_string += arg
		if i != len(args)-1 {
			cmd_string += " "
		}
	}
	log("cmd: ", cmd)
	log("args:", args)
	log("cmd_string", cmd_string)

	command = exec.Command(cmd, args...)

	command.Stdout = &buf_stdout
	command.Stderr = &buf_stderr
	if err = command.Run(); err != nil {
		ret = -1
	} else {
		ret = 0
	}
	stdout = buf_stdout.String()
	stderr = buf_stderr.String()

	return
}

//export Execv1
func Execv1(cmd string, args string, shell bool) (ret int, stdout string, stderr string) {
	return Execv(cmd, SliceIt(args), shell)
}
