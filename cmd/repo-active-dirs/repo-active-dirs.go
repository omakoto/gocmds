///bin/true; exec /usr/bin/env go run "$0" "$@"

// repo-active-dirs
// Find repo projects with one or more local branches and print them.

package main

import (
	"bytes"
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/repo"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"
)

func realMain() int {
	manifest, root := repo.MustLoadManifest()

	ways := 8 // Number of parallel goroutines.

	mu := sync.Mutex{} // Protect list.
	list := make([]string, 0, len(manifest.Projects))

	ch := make(chan string, ways*2)

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
				list = append(list, dir)
				mu.Unlock()
			}
			wg.Done()
		}()
	}

	// Read all projects
	for _, p := range manifest.Projects {
		// fmt.Print(filepath.Join(root, p.Path))
		// fmt.Print("\n")
		ch <- p.Path
	}
	for i := 0; i < ways; i++ {
		ch <- ""
	}
	wg.Wait()

	// Sort and print
	sort.Strings(list)

	// Print all.
	for _, l := range list {
		fmt.Print(l)
		fmt.Print("\n")
	}

	return 0
}

func main() {
	common.RunAndExit(realMain)
}
