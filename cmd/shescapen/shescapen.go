///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/cmd/shescapecommon"
	"os"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	shescapecommon.ShescapeNoNewline(os.Args)
	return 0
}
