package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/omakoto/go-common/src/common"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {

	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: fsub DECIMAL1 DECIMAL2 (prints DECIMAL1 - DECIMAL2)\n")
		return 1
	}
	v1, err := strconv.ParseFloat(args[0], 32)
	common.Checke(err)

	v2, err := strconv.ParseFloat(args[1], 32)
	common.Checke(err)

	fmt.Printf("%.6f\n", (v1 - v2))
	return 0
}
