package repo

import (
	"path/filepath"
	"github.com/omakoto/gaze/src/common"
	"os"
)

const EnvBuildTop = "ANDROID_BUILD_TOP"

func MustFindRepoTop(path string) string {
	atop := common.MustGetenv(EnvBuildTop)

	path, err := filepath.Abs(path)
	common.Check(err, "Abs() failed")
	for {
		s, err := os.Stat(filepath.Join(path, ".repo"))
		if err == nil && s.IsDir() {
			if atop != path {
				common.Fatal("Not in " + EnvBuildTop)
			}
			return path
		}
		if path == "/" {
			common.Fatal("Repo top directory not found")
		}
		path = filepath.Dir(path)
	}
}
