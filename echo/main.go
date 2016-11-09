package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("args: %v", strings.Join(os.Args, " "))
}
