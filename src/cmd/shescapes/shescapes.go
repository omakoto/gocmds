///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/shescapecommon"
	"os"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	shescapecommon.ShescapeStdin(os.Args)
	return 0
}
