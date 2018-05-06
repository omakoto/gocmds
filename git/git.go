package git

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"os"
	"path/filepath"
)

func FindGitTop(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	for {
		s, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil && s.IsDir() {
			return path, nil
		}
		if path == "/" {
			return "", fmt.Errorf("Git top directory not found")
		}
		path = filepath.Dir(path)
	}
}

func MustFindGitTop(path string) string {
	root, err := FindGitTop(path)
	if err != nil {
		common.Fatalf("%s", err)
	}
	return root
}
