///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio"
	"github.com/omakoto/gocmds/src/cmd/shescapecommon"
	"os"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	shescapecommon.ShescapeNoNewline(os.Args[1:])
	textio.BufferedStdout.WriteByte('\n')
	return 0
}
