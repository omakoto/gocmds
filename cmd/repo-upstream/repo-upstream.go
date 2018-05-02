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
	"flag"
)

var (
	quiet = flag.Bool("q", false, "Quiet mode")
	origin = flag.Bool("o", false, "Only show the default source")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	flag.Parse()
	common.Quiet = *quiet

	manifest, root := repo.MustLoadManifest()

	path, err := os.Getwd()
	common.Check(err, "Getwd() failed")
	if len(flag.Args()) > 0 {
		path = flag.Args()[0]

		s, err := os.Stat(path)
		common.Checkf(err, "Directory \"%s\" doesn't exist", path)
		if !s.IsDir() {
			common.Fatalf("Not a directory: \"%s\"", path)
		}
		path, err = filepath.Abs(path)
		common.Checkf(err, "Abs() failed")
	}
	common.Debugf("root=%s\n", root)
	common.Debugf("path=%s\n", path)
	common.Dump("manifest=", manifest)

	remote := manifest.Default.Remote
	revision := manifest.Default.Revision

	if !*origin {
		gitTop, err := git.FindGitTop(path)
		if err == nil {
			common.Debugf("git-top=%s\n", gitTop)

			for _, p := range manifest.Projects {
				if gitTop == filepath.Join(root, p.Path) {
					if p.Remote != "" {
						remote = p.Remote
					}
					if p.Revision != "" {
						revision = p.Revision
					}
					break
				}
			}
		}
	}
	// We still want to show it even at the top directory where we're not in any project.
	fmt.Print(remote, "/", revision, "\n")
	return 0
}
