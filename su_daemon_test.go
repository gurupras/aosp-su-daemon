package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/flynn-archive/go-shlex"
	"github.com/gurupras/aosp_su_daemon"
	"github.com/gurupras/testing_base"
)

func TestWrite(t *testing.T) {
	var success bool = true
	var args []string
	var err error
	result := testing_base.InitResult("TestWrite-1")
	if args, err = (shlex.Split("su_daemon write 'hello world' test")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret := main.Main(args); ret < 0 {
			success = false
		}
	}

	// Out write function cannot create new files. So test this out.
	result = testing_base.InitResult("TestWrite-2")
	if args, err = (shlex.Split("su_daemon write 'hello world' /proc/a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret := main.Main(args); ret >= 0 {
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
	result := testing_base.InitResult("TestExec-1")
	if args, err = (shlex.Split("su_daemon -v exec ls -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret := main.Main(args); ret < 0 {
			success = false
		}
	}
	result = testing_base.InitResult("TestExec-2")
	if args, err = (shlex.Split("su_daemon exec programmustnotexist -- -l -i -s -a")); err != nil {
		fmt.Fprintln(os.Stderr, err)
		success = false
		goto out
	} else {
		if ret := main.Main(args); ret >= 0 {
			success = false
		}
	}
out:
	testing_base.HandleResult(t, success, result)
}
