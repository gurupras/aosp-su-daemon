package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/google/shlex"
)

var (
	UnixSocketPath = "/dev/socket/su_daemon"
	LogPath        = "/dev/kmsg"
	socket         net.Listener
	log_buf        *bufio.Writer
)

const (
	TAG = "su_daemon"
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

var (
	app     = kingpin.New("su-daemon", "An AOSP daemon with su privileges")
	verbose = app.Flag("verbose", "Enable verbose output").Short('v').Default("false").Bool()
	output  = app.Flag("output", "Print output").Short('o').Default("false").Bool()

	write_cmd  = app.Command("write", "Write data to file")
	write_data = write_cmd.Arg("data", "Data to write").Required().String()
	write_file = write_cmd.Arg("file", "File to write to").Required().String()

	exec_cmd      = app.Command("exec", "Execute a program")
	exec_cmd_bin  = exec_cmd.Arg("binary", "Binary to execute").Required().String()
	exec_cmd_args = NArgs(exec_cmd.Arg("args", "Arguments to binary"))
)

func debug(a ...interface{}) {
	if *verbose {
		fmt.Println(a...)
	}
}

func SliceIt(args string) (ret []string) {
	ret, _ = shlex.Split(args)
	return
}

func init_fn() {
	var err error
	if socket != nil {
		goto log
	}
	socket, err = net.Listen("unix", UnixSocketPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to bind to socket:", err)
		os.Exit(-1)
	}
	if err = os.Chmod(UnixSocketPath, 0666); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to chmod the socket:", err)
		os.Exit(-1)
	}
log:
	if log_buf != nil {
		return
	}
	log_file, err := os.OpenFile(LogPath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to open:", LogPath))
		os.Exit(-1)
	}
	log_buf = bufio.NewWriter(log_file)
}

//export Main1
func Main1(command string) (ret int, stdout string, stderr string) {
	args := SliceIt(command)
	return Main(args)
}

func Main(args []string) (ret int, stdout string, stderr string) {
	init_fn()
	switch kingpin.MustParse(app.Parse(args[1:])) {
	case write_cmd.FullCommand():
		debug("write %s %s", write_data, write_file)
		ret = Write(*write_data, *write_file)
	case exec_cmd.FullCommand():
		debug("exec:", *exec_cmd_bin, *exec_cmd_args)
		ret, stdout, stderr = Execv(*exec_cmd_bin, *exec_cmd_args)
	}
	if *output {
		log(stdout)
		log(stderr)
	}
	return
}

func log(msg ...interface{}) {
	log_buf.Write([]byte(fmt.Sprintf("%v: %v\n", TAG, msg)))
	log_buf.Flush()
}

func process(fd net.Conn) {
	buf := make([]byte, 512)
	nr, err := fd.Read(buf)
	if err != nil {
		return
	}

	cmd_bytes := buf[0:nr]
	cmd := string(cmd_bytes)
	log("Command: %s", cmd)
	Main1(cmd)
}

func server() {
	for {
		fd, err := socket.Accept()
		log("Received connection")
		if err != nil {
			log("Failed to accept socket")
			continue
		}
		process(fd)
		fd.Close()
	}
}

func main() {
	init_fn()
	server()
}
