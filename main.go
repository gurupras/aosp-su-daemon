package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/google/shlex"
)

var (
	ShellPath      = "/system/bin/sh"
	UnixSocketPath = "/dev/socket/su_daemon"
	LogPath        = "/dev/kmsg"
	socket         *net.UnixListener
	LogBuf         *bufio.Writer
)

const (
	TAG = "su_daemon"
)

var (
	app           *kingpin.Application
	verbose       *bool
	output        *bool
	write_cmd     *kingpin.CmdClause
	write_data    *string
	write_file    *string
	exec_cmd      *kingpin.CmdClause
	exec_cmd_bin  *string
	exec_cmd_args *[]string
	exec_shell    *bool
)

type nargs []string

func (n *nargs) Set(v string) error {
	*n = append(*n, v)
	return nil
}

func (n *nargs) String() string {
	return ""
}

func (n *nargs) IsCumulative() bool {
	return true
}

func NArgs(s kingpin.Settings) (target *[]string) {
	target = new([]string)
	s.SetValue((*nargs)(target))
	return
}

func init_kingpin() {
	app = kingpin.New("su-daemon", "An AOSP daemon with su privileges")
	verbose = app.Flag("verbose", "Enable verbose output").Short('v').Default("false").Bool()
	output = app.Flag("output", "Print output").Short('o').Default("false").Bool()

	write_cmd = app.Command("write", "Write data to file")
	write_data = write_cmd.Arg("data", "Data to write").Required().String()
	write_file = write_cmd.Arg("file", "File to write to").Required().String()

	exec_cmd = app.Command("exec", "Execute a program")
	exec_cmd_bin = exec_cmd.Arg("binary", "Binary to execute").Required().String()
	exec_cmd_args = NArgs(exec_cmd.Arg("args", "Arguments to binary"))
	exec_shell = app.Flag("shell", "Run command in a shell").Short('s').Default("false").Bool()
}

func debug(a ...interface{}) {
	if *verbose {
		fmt.Println(a...)
	}
}

func log(msg ...interface{}) {
	LogBuf.Write([]byte(fmt.Sprintf("%v: %v\n", TAG, msg)))
	LogBuf.Flush()
}

func SliceIt(args string) (ret []string) {
	ret, _ = shlex.Split(args)
	return
}

func init_fn() {
	var err error

	// Initialize kingpin
	init_kingpin()

	if socket != nil {
		goto log
	}
	socket, err = net.ListenUnix("unix", &net.UnixAddr{UnixSocketPath, "unix"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to bind to socket:", err)
		os.Exit(-1)
	}
	if err = os.Chmod(UnixSocketPath, 0666); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to chmod the socket:", err)
		os.Exit(-1)
	}
log:
	if LogBuf != nil {
		return
	}
	log_file, err := os.OpenFile(LogPath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to open:", LogPath))
		os.Exit(-1)
	}
	LogBuf = bufio.NewWriter(log_file)
}

//export Main1
func Main1(command string) (ret int, stdout string, stderr string) {
	args := SliceIt(command)
	return Main(args)
}

func Main(args []string) (ret int, stdout string, stderr string) {
	init_fn()
	cmd, err := app.Parse(args[1:])
	if err != nil {
		ret = -1
		stderr = err.Error()
		goto out
	}
	switch kingpin.MustParse(cmd, err) {
	case write_cmd.FullCommand():
		debug("write:", *write_data, *write_file)
		ret = Write(*write_data, *write_file)
	case exec_cmd.FullCommand():
		debug("exec:", *exec_cmd_bin, *exec_cmd_args)
		ret, stdout, stderr = Execv(*exec_cmd_bin, *exec_cmd_args, *exec_shell)
	}
	if *output {
		log(stdout)
		log(stderr)
	}
out:
	return
}

func process(fd net.Conn) {
	buf := make([]byte, 512)
	nr, err := fd.Read(buf)
	if err != nil {
		return
	}

	cmd_bytes := buf[0:nr]
	cmd := strings.TrimSpace(string(cmd_bytes))
	log("cmd_bytes:", cmd)
	ret, stdout, stderr := Main1(cmd)
	var retbuf bytes.Buffer

	var total_len uint64

	total_len = 0 +
		8 /* total */ +
		4 /* ret */ +
		8 /* len(stdout) */ +
		uint64(len(stdout)) /* stdout */ +
		8 /* len(stderr) */ +
		uint64(len(stderr)) /* stderr */

	fmt.Println("total_len:", total_len)
	if err = binary.Write(&retbuf, binary.LittleEndian, total_len); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}
	if err = binary.Write(&retbuf, binary.LittleEndian, int32(ret)); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}

	if err = binary.Write(&retbuf, binary.LittleEndian, uint64(len(stdout))); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}
	if _, err = retbuf.WriteString(stdout); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}

	if err = binary.Write(&retbuf, binary.LittleEndian, uint64(len(stderr))); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}
	if _, err = retbuf.WriteString(stderr); err != nil {
		fmt.Println("Failed to write to bytes.Buffer", err)
	}

	debug("retval:", ret)
	debug("stdout:", stdout)
	debug("stderr:", stderr)

	n, err := fd.Write(retbuf.Bytes())
	_ = n
	debug("Wrote:", n)
	fd.Close()
}

func server() {
	for {
		fd, err := socket.AcceptUnix()
		log("Received connection")
		if err != nil {
			log("Failed to accept socket")
			continue
		}
		go process(fd)
	}
}

func main() {
	init_fn()
	server()
}
