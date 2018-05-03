package main

import (
	"github.com/omakoto/gaze/src/common"
	"github.com/omakoto/gocmds/repo"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	common.Quiet = true

	repo.MustLoadManifest()
	return 0
}
