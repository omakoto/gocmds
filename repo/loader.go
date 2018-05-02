package repo

import (
	"github.com/omakoto/gaze/src/common"
	"io/ioutil"
	"path/filepath"
	"encoding/xml"
)

type Remote struct {
	Name   string `xml:"name,attr"`
	Fetch  string `xml:"fetch,attr"`
	Review string `xml:"review,attr"`
}

type Default struct {
	Revision   string `xml:"revision,attr"`
	Remote     string `xml:"remote,attr"`
	DestBranch string `xml:"dest-branch,attr"`
}

type Project struct {
	Path       string `xml:"path,attr"`
	Name       string `xml:"name,attr"`
	Revision   string `xml:"revision,attr"`
	Remote     string `xml:"remote,attr"`
	DestBranch string `xml:"dest-branch,attr"`
}

type Manifest struct {
	Default  Default   `xml:"default"`
	Remotes  []Remote  `xml:"remote"`
	Projects []Project `xml:"project"`
}

func MustLoadManifest() (manifest *Manifest, root string) {
	//root = os.Getenv("ANDROID_BUILD_TOP")
	//common.Debugf("ANDROID_BUILD_TOP=%s\n", root)

	root = MustFindRepoTop(".")

	// Find all projects.
	manifestRaw, err := ioutil.ReadFile(filepath.Join(root, ".repo/manifest.xml"))
	common.Check(err, "Failed to read manifest.xml")

	manifest = &Manifest{}
	common.Check(xml.Unmarshal(manifestRaw, manifest), "Failed to parse manifest.xml")
	return
}
