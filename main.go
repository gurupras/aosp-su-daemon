package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
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

func main() {
	Main(os.Args)
}
