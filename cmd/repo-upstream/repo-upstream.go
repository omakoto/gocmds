///bin/true; exec /usr/bin/env go run "$0" "$@"

// repo-upstream
// Print the repo upstream branch name for the current path, or a given path.

package main

import (
	"github.com/omakoto/gaze/src/common"
	"github.com/omakoto/gocmds/repo"
	"os"
	"path/filepath"
	"github.com/omakoto/gocmds/git"
	"fmt"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	manifest, root := repo.MustLoadManifest()

	path, err := os.Getwd()
	common.Check(err, "Getwd() failed")
	if len(os.Args) > 1 {
		path = os.Args[1]

		s, err := os.Stat(path)
		common.Checkf(err, "Directory \"%s\" doesn't exist", path)
		if !s.IsDir() {
			common.Fatalf("Not a directory: \"%s\"", path)
		}
		path, err = filepath.Abs(path)
		common.Checkf(err, "Abs() failed")
	}
	gitTop := git.MustFindGitTop(path)
	common.Debugf("root=%s\n", root)
	common.Debugf("path=%s\n", path)
	common.Debugf("git-top=%s\n", gitTop)
	common.Dump("manifest=", manifest)

	remote := manifest.Default.Remote
	revision := manifest.Default.Revision
	for _, p := range manifest.Projects {
		if gitTop == filepath.Join(root, p.Path) {
			if p.Remote != "" {
				remote = p.Remote
			}
			if p.Revision != "" {
				revision = p.Revision
			}
			fmt.Print(remote, "/", revision, "\n")
			return 0
		}
	}
	// common.Fatal("Upstream not found")
	return 1
}
