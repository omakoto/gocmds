///bin/true; exec /usr/bin/env go run "$0" "$@"

// git-top
// Prints the git top directory.

package main

import (
	"github.com/omakoto/gaze/src/common"
	"fmt"
	"github.com/omakoto/gocmds/git"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	fmt.Print(git.MustFindGitTop("."), "\n")
	return 0
}