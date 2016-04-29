package main

import (
	"bufio"
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
	var err error
	kmsg_file, err := os.OpenFile("/dev/kmsg", os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open /dev/kmsg")
		os.Exit(-1)
	}
	kmsg = bufio.NewWriter(kmsg_file)
}

func log(msg ...interface{}) {
	kmsg.Write([]byte(fmt.Sprintf("%v: %v\n", TAG, msg)))
	kmsg.Flush()
}

func main() {
	var err error
	c, err := net.Dial("unix", UNIX_SOCKET)
	if err != nil {
		log("Failed to dial socket")
	}

	_, err = c.Write([]byte("su -vo write 125 /sys/class/leds/led_flash_torch/brightness"))
	if err != nil {
		log("Failed to turn on flashlight: %v", err)
	}
}
