package main

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/repo"
	"github.com/pborman/getopt/v2"
)

var (
	verbose = getopt.BoolLong("verbose", 'v', "Show error message.")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()
	common.Quiet = !*verbose

	repo.MustLoadManifest()

	return 0
}
