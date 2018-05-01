///bin/true; exec /usr/bin/env go run "$0" "$@"

// repo-active-dirs
// Find repo projects with one or more local branches and print them.

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/omakoto/gaze/src/common"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"bytes"
)

type Project struct {
	Path string `xml:"path,attr"`
}

type Manifest struct {
	Projects []Project `xml:"project"`
}

func realMain() int {
	// Find the top diir.
	root := os.Getenv("ANDROID_BUILD_TOP")
	common.Debugf("ANDROID_BUILD_TOP=%s", root)

	// Find all projects.
	manifest, err := ioutil.ReadFile(filepath.Join(root, ".repo/manifest.xml"))
	common.Check(err, "Failed to read manifest.xml")

	var man Manifest
	xml.Unmarshal(manifest, &man)

	mu := sync.Mutex{} // Control output.
	ways := 8 // Number of parallel goroutines.

	ch := make(chan string, ways * 2)
	
	wg := sync.WaitGroup{}
	wg.Add(ways)

	for i := 0; i < ways; i++ {
		go func() {
			for path := range ch {
				if path == "" {
					break
				}
				// Read .git/config and find "[branch ".
				dir := filepath.Join(root, path)
				config := filepath.Join(dir, ".git/config")
				data, err := ioutil.ReadFile(config)
				if err != nil {
					continue // Find not found or not readable.
				}
				if bytes.Index(data, []byte("[branch \"")) < 0 {
					continue // No local branches.
				}

				mu.Lock()
				fmt.Print(dir)
				fmt.Print("\n")
				mu.Unlock()
			}
			wg.Done()
		}()
	}

	// Read all projects
	for _, p := range man.Projects {
		// fmt.Print(filepath.Join(root, p.Path))
		// fmt.Print("\n")
		ch <- p.Path
	}
	for i := 0; i < ways; i++ {
		ch <- ""
	}
	wg.Wait()

	return 0
}

func main() {
	common.RunAndExit(realMain)
}
