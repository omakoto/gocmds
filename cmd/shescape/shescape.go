///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/cmd/shescapecommon"
	"github.com/pborman/getopt/v2"
)

var (
	noNewline = getopt.BoolLong("no-newline", 'n', "Don't print newline.")
	fromStdin = getopt.BoolLong("stdin", 's', "Read input from stdin.")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	if !*fromStdin {
		shescapecommon.ShescapeNoNewline(getopt.Args())
		if !*noNewline {
			fmt.Println()
		}
	} else {
		shescapecommon.ShescapeStdin(getopt.Args())
	}
	return 0
}
