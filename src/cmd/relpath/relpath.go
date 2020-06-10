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

var (
	basePath = getopt.StringLong("base-path", 'p', "Base directory")
)

func init() {
	getopt.SetUsage(usage)
}

func usage() {
	os.Stderr.WriteString(`
relpath: Convert PATHs to relative paths

Usage: relpath [-p BASE-DIR] PATH [...]

`)
	getopt.CommandLine.PrintOptions(os.Stderr)
}

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	base := *basePath
	if base == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		base = cwd
	}

	ret := 0
	for _, file := range getopt.Args() {
		rel, err := filepath.Rel(base, utils.HomeExpanded(file))
		if err != nil {
			common.Warnf("%s\n", err)
			ret = 1
			continue
		}
		fmt.Print(rel, "\n")
	}
	return ret
}
