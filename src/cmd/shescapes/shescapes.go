package main

import (
	"os"

	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/shescapecommon"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	shescapecommon.ShescapeStdin(os.Args[1:])
	return 0
}
