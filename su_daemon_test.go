package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/flynn-archive/go-shlex"
	"github.com/gurupras/gocommons"
)

func init() {
	var f *os.File
	var err error

	UnixSocketPath = "/tmp/test.sock"
	// Delete socket if it exists
	if _, err = os.Stat(UnixSocketPath); err == nil {
		//fmt.Println("Deleting existing socket file")
		os.Remove(UnixSocketPath)
	}
	//fmt.Println("Socket:", UnixSocketPath)

	// Create log file if it doesn't exist
	LogPath = "/tmp/log"
	if f, err = os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Could not create log file: %s", LogPath))
		os.Exit(-1)
	}
	LogBuf = bufio.NewWriter(f)
	init_fn()
	gocommons.ShellPath = "/bin/bash"
}

func TestWrite(t *testing.T) {
	var success bool = true
	var args []string
	var err error
	var result string

	// Create a file that we will write to
	var f *os.File
	if f, err = os.OpenFile("test", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Could not create file to test")
		success = false
		goto out
	}
	f.Close()
	result = gocommons.InitResult("TestWrite-1")
	if args, err = (shlex.Split("su_daemon -v write 'hello world' test")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret < 0 {
			success = false
		}
	}
	os.Remove("test")
	gocommons.HandleResult(t, success, result)

	// Out write function cannot create new files. So test this out.
	result = gocommons.InitResult("TestWrite-2")
	if args, err = (shlex.Split("su_daemon -v write 'hello world' /proc/a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret >= 0 {
			success = false
		}
	}
out:
	gocommons.HandleResult(t, success, result)
}

func TestExec(t *testing.T) {
	var success bool = true
	var args []string
	var err error
	var result string

	result = gocommons.InitResult("TestExec-1")
	if args, err = (shlex.Split("su_daemon -vo exec ls -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret < 0 {
			success = false
		}
	}
	gocommons.HandleResult(t, success, result)

	result = gocommons.InitResult("TestExec-2")
	if args, err = (shlex.Split("su_daemon -vo exec programmustnotexist -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret >= 0 {
			success = false
		}
	}
out:
	gocommons.HandleResult(t, success, result)
}

func TestExecShell(t *testing.T) {
	var success bool = true
	var args []string
	var err error
	var result string

	_ = "breakpoint"
	result = gocommons.InitResult("TestExecShell-1")
	if args, err = (shlex.Split("su_daemon -vso exec ls -- -l -i -s -a")); err != nil {
		fmt.Println(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret < 0 {
			success = false
		}
	}
	gocommons.HandleResult(t, success, result)

	result = gocommons.InitResult("TestExecShell-2")
	if args, err = (shlex.Split("su_daemon -vso exec programmustnotexist -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret >= 0 {
			success = false
		}
	}
out:
	gocommons.HandleResult(t, success, result)
}
