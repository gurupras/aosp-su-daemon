package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gurupras/gocommons"
)

//export Write
func Write(data, file string) (ret int, stdout string, stderr string) {
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
		stderr = fmt.Sprintf("Failed to write to '%v': %v", file, err)
		fmt.Fprintln(os.Stderr, stderr)
		ret = -1
		return
	}
	return
}

//export Read
func Read(file string) (ret int, stdout string, stderr string) {
	var f_raw *os.File
	var err error

	if f_raw, err = os.OpenFile(file, os.O_RDONLY, 0); err != nil {
		stderr = fmt.Sprintf("Failed to open file '%v': %v", file, err)
		fmt.Fprintln(os.Stderr, stderr)
		ret = -1
		return
	}
	defer f_raw.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, f_raw); err != nil {
		stderr = fmt.Sprintf("Failed to read bytes from '%v': %v", file, err)
		fmt.Fprintln(os.Stderr, stderr)
		ret = -1
		return
	}
	stdout = strings.TrimSpace(buf.String())
	return
}

//export RunUiAutomator
func RunUiAutomator(jarPath, jarMethod, extras string) (err error) {
	cmd := "/system/bin/uiautomator"
	args := fmt.Sprintf("runtest %s %s -c %s", jarPath, extras, jarMethod)
	log(fmt.Sprintf("%s %s", cmd, args))
	ret, stdout, stderr := gocommons.Execv1(cmd, args, true)
	_ = ret
	_ = stdout
	_ = stderr
	return nil
}
