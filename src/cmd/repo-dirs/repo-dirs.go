///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/repo"
	"github.com/ungerik/go-dry"
	"os"
	"path"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	manifest, root := repo.MustLoadManifest()
	mirror := repo.IsMirror()

	dry.Nop(mirror)

	status := 0

	for _, p := range manifest.Projects {
		var dir string
		if !mirror {
			dir = path.Join(root, p.Path)
		} else {
			dir = path.Join(root, p.Name) + ".git"
		}
		if dry.FileIsDir(dir) {
			fmt.Print(dir, "\n")
		} else {
			fmt.Fprint(os.Stderr, dir, " doesn't exist\n")
			status = 1
		}
	}

	return status
}
