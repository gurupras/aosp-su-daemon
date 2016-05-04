package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/flynn-archive/go-shlex"
	"github.com/gurupras/testing_base"
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
	defer f.Close()
	init_fn()
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
	result = testing_base.InitResult("TestWrite-1")
	if args, err = (shlex.Split("su_daemon write 'hello world' test")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret < 0 {
			success = false
		}
	}
	os.Remove("test")
	testing_base.HandleResult(t, success, result)

	// Out write function cannot create new files. So test this out.
	result = testing_base.InitResult("TestWrite-2")
	if args, err = (shlex.Split("su_daemon write 'hello world' /proc/a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret >= 0 {
			success = false
		}
	}
out:
	testing_base.HandleResult(t, success, result)
}

func TestExec(t *testing.T) {
	var success bool = true
	var args []string
	var err error
	var result string

	result = testing_base.InitResult("TestExec-1")
	if args, err = (shlex.Split("su_daemon -v exec ls -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret < 0 {
			success = false
		}
	}
	testing_base.HandleResult(t, success, result)

	result = testing_base.InitResult("TestExec-2")
	if args, err = (shlex.Split("su_daemon exec programmustnotexist -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret, _, _ := Main(args); ret >= 0 {
			success = false
		}
	}
out:
	testing_base.HandleResult(t, success, result)
}
