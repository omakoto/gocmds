package repo

import (
	"path/filepath"
	"github.com/omakoto/gaze/src/common"
	"os"
)

func MustFindRepoTop(path string) string {
	path, err := filepath.Abs(path)
	common.Check(err, "Abs() failed")
	for {
		s, err := os.Stat(filepath.Join(path, ".repo"))
		if err == nil && s.IsDir() {
			return path
		}
		if path == "/" {
			common.Fatal("Repo top directory not found")
		}
		path = filepath.Dir(path)
	}
}
