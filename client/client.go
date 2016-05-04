package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

var (
	UNIX_SOCKET = "/dev/socket/su_daemon"
	kmsg        *bufio.Writer
)

const (
	TAG = "su_client"
)

func init() {
	kmsg = bufio.NewWriter(os.Stdout)
}

func log(msg ...interface{}) {
	kmsg.Write([]byte(fmt.Sprintf("%v: %v\n", TAG, msg)))
	kmsg.Flush()
}

func main() {
	var c net.Conn
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Connecting to su_daemon...")
		c, err = net.DialUnix("unix", nil, &net.UnixAddr{UNIX_SOCKET, "unix"})
		if err != nil {
			log("Failed to dial socket")
		}
		fmt.Printf("Connected\n")

		cmd, _ := reader.ReadString('\n')
		_, err = c.Write([]byte(cmd))
		if err != nil {
			log("Failed to run:", cmd, err)
		}

		total_len_b := make([]byte, 8)
		ret_b := make([]byte, 4)
		stdout_len_b := make([]byte, 8)
		stderr_len_b := make([]byte, 8)

		if _, err = c.Read(total_len_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive response from server:", err)
		}
		total_len := binary.LittleEndian.Uint64(total_len_b)
		fmt.Println("total_len:", total_len)

		if _, err = c.Read(ret_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive return code from server:", err)
		}
		ret := int32(binary.LittleEndian.Uint32(ret_b))
		fmt.Println("ret:", ret)

		if _, err = c.Read(stdout_len_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive stdout length from server:", err)
		}
		stdout_len := binary.LittleEndian.Uint64(stdout_len_b)
		fmt.Println("stdout_len:", stdout_len)
		stdout_b := make([]byte, stdout_len)
		if _, err = c.Read(stdout_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive stdout from server:", err)
		}
		stdout := string(stdout_b)
		fmt.Println("stdout:", stdout)

		if _, err = c.Read(stderr_len_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive stderr length from server:", err)
		}
		stderr_len := binary.LittleEndian.Uint64(stderr_len_b)
		fmt.Println("stderr_len:", stderr_len)
		stderr_b := make([]byte, stderr_len)
		if _, err = c.Read(stderr_b); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to receive stderr from server:", err)
		}
		stderr := string(stderr_b)
		fmt.Println("stderr:", stderr)

		_ = ret
		fmt.Fprintln(os.Stdout, stdout)
		fmt.Fprintln(os.Stderr, stderr)
	}
}
