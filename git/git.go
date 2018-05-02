package git

import (
	"os"
	"path/filepath"
	"github.com/omakoto/gaze/src/common"
)

func MustFindGitTop(path string) string {
	path, err := filepath.Abs(path)
	common.Check(err, "Abs() failed")
	for {
		s, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil && s.IsDir() {
			return path
		}
		if path == "/" {
			common.Fatal("Git top directory not found")
		}
		path = filepath.Dir(path)
	}
}
