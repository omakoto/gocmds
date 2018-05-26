package repo

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/fileutils"
	"os"
	"path/filepath"
)

const EnvBuildTop = "ANDROID_BUILD_TOP"

func FindRepoTop(path string) (string, error) {
	atop := os.Getenv(EnvBuildTop)

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	for {
		s, err := os.Stat(filepath.Join(path, ".repo"))
		if err == nil && s.IsDir() {
			if atop != "" && !fileutils.SamePath(atop, path) {
				return "", fmt.Errorf("not in $%s", EnvBuildTop)
			}
			return path, nil
		}
		if path == "/" {
			return "", fmt.Errorf("repo top directory not found")
		}
		path = filepath.Dir(path)
	}
}

func MustFindRepoTop(path string) string {
	ret, err := FindRepoTop(path)
	common.Check(err, "Not in repo")
	return ret
}
