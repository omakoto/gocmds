///bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/utils"
	"github.com/pborman/getopt/v2"
	"os"
	"path/filepath"
)

func init() {
	getopt.SetUsage(usage)
}

func usage() {
	os.Stderr.WriteString(`
abspath: Convert PATHs to absolute paths

Usage: abspath PATH [...]

`)
	getopt.CommandLine.PrintOptions(os.Stderr)
}

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	ret := 0
	for _, file := range getopt.Args() {
		abs, err := filepath.Abs(utils.HomeExpanded(file))
		if err != nil {
			common.Warnf("%s\n", err)
			ret = 1
			continue
		}
		fmt.Print(abs)
		fmt.Print("\n")
	}
	return ret
}
