package repo

import (
	"encoding/xml"
	"github.com/omakoto/go-common/src/common"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
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

type Include struct {
	Name string `xml:"name,attr"`
}

type Manifest struct {
	Default  Default   `xml:"default"`
	Remotes  []Remote  `xml:"remote"`
	Projects []Project `xml:"project"`
	Include  Include   `xml:"include"`
}

func LoadManifestFromFile(file string) (manifest *Manifest, root string, err error) {
	root, err = FindRepoTop(".")
	if err != nil {
		return nil, "", err
	}

	// Find all projects.
	manifestRaw, err := ioutil.ReadFile(filepath.Join(root, file))
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to read manifest.xml")
	}

	manifest = &Manifest{}
	err = xml.Unmarshal(manifestRaw, manifest)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse manifest.xml")
	}

	return
}

func LoadManifest() (manifest *Manifest, root string, err error) {
	manifest, root, err = LoadManifestFromFile(".repo/manifest.xml")
	if err != nil {
		return
	}
	if len(manifest.Projects) == 0 && manifest.Include.Name != "" {
		manifest, root, err = LoadManifestFromFile(".repo/manifests/" + manifest.Include.Name)
	}
	return
}

func MustLoadManifest() (manifest *Manifest, root string) {
	manifest, root, err := LoadManifest()
	common.Check(err, "Not in repo")
	return
}
