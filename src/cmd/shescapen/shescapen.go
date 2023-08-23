// bin/true; exec /usr/bin/env go run "$0" "$@"
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
	shescapecommon.ShescapeNoNewline(os.Args[1:])
	return 0
}
