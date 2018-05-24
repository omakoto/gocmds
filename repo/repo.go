package repo

import (
	"os/exec"
	"path"
	"strings"
)

func IsMirror() bool {
	root := MustFindRepoTop(".")

	out, err := exec.Command("git", "-C", path.Join(root, ".repo/manifests.git/"), "config", "--local", "--get", "repo.mirror").Output()
	if err != nil {
		return false
	}
	// common.Check(err, "git config failed.")

	return strings.TrimSpace(string(out)) == "true"
}
