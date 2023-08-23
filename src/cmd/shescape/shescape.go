package main

import (
	"os"

	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio"
	"github.com/omakoto/gocmds/src/cmd/shescapecommon"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	shescapecommon.ShescapeNoNewline(os.Args[1:])
	textio.BufferedStdout.WriteByte('\n')
	return 0
}
