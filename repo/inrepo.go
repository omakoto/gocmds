package repo

import (
	"path/filepath"
	"github.com/omakoto/go-common/src/common"
	"os"
	"github.com/omakoto/go-common/src/fileutils"
)

const EnvBuildTop = "ANDROID_BUILD_TOP"

func MustFindRepoTop(path string) string {
	atop := common.MustGetenv(EnvBuildTop)

	path, err := filepath.Abs(path)
	common.Check(err, "Abs() failed")
	for {
		common.Debugf("path=%s\n", path)
		s, err := os.Stat(filepath.Join(path, ".repo"))
		if err == nil && s.IsDir() {
			if !fileutils.SamePath(atop, path) {
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
