package main

import (
	"bufio"
	"fmt"
	"github.com/omakoto/go-common/src/must"
	"os"
	"strings"
)

func FindDirs(out *bufio.Writer, dir string) {
	d := dir
	if !strings.HasSuffix(d, "/") {
		d = dir + "/"
	}
	findDirs(out, d)
}

func findDirs(out *bufio.Writer, dir string) {
	entries := must.Must2(os.ReadDir(dir))

	for _, entry := range entries {
		if entry.IsDir() {
			d := dir + entry.Name() + "/"
			out.WriteString(d)
			out.WriteByte('\n')
			FindDirs(out, d)
		}
	}
}

func main() {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	dirs := []string{"."}
	if len(os.Args) > 1 {
		dirs = os.Args[1:]
	}

	for _, dir := range dirs {
		fmt.Printf("# %s\n", dir)
		FindDirs(out, dir)
	}
}
