///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

// Show sizes of all files under given directories.

import (
	"fmt"
	"os"
	"path/filepath"
)

func showFiles(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			fmt.Printf("%d\t%s\n", info.Size(), path)
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "filesize: %s\n", err)
				return nil
			}
			fmt.Printf("%d\t%s -> %s\n", info.Size(), path, target)
		}
		return nil
	})
}

func main() {
	for _, dir := range os.Args[1:] {
		showFiles(dir)
	}
}
