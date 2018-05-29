///bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/ffind/printer"
	"github.com/pborman/getopt/v2"
	"os"
	"path"
	"runtime"
	"sort"
	"sync"
)

const maxPara = 32

func defaultPara() int {
	v := runtime.NumCPU() * 4
	if v < maxPara {
		return v
	}
	return maxPara
}

var (
	showFiles = getopt.BoolLong("file", 'f', "Print files only")
	showDirs  = getopt.BoolLong("dir", 'd', "Print directories only")

	quiet = getopt.BoolLong("quiet", 'q', "Don't show warnings")

	para = getopt.IntLong("para", 'j', defaultPara(), "Number of goroutines")

	ch = make(chan string, 1024*1024)

	numBacklog = 0
	cond       = sync.NewCond(&sync.Mutex{})
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()
	defer printer.Flush()

	common.Quiet = *quiet

	if !(*showFiles || *showDirs) {
		*showFiles = true
		*showDirs = true
	}

	common.Debugf("-j=%d\n", *para)

	for i := 0; i < *para; i++ {
		go func() {
			for {
				dir := <-ch
				//common.Debugf("Pop:  %s\n", dir)
				doFindDir(dir)

				cond.L.Lock()
				numBacklog--
				//common.Debugf("Done: %s [%d]\n", dir, numBacklog)
				if numBacklog <= 0 {
					cond.Signal()
				}
				cond.L.Unlock()
			}
		}()
	}

	if len(getopt.Args()) == 0 {
		dir, err := os.Getwd()
		common.Check(err, "Getwd failed")

		findDir(dir)
	} else {
		for _, dir := range getopt.Args() {
			findDir(dir)
		}
	}

	cond.L.Lock()
	if numBacklog > 0 {
		cond.Wait()
	}
	cond.L.Unlock()

	return 0
}

func findDir(dir string) {
	cond.L.Lock()
	numBacklog++
	//common.Debugf("Push: %s [%d]\n", dir, numBacklog)
	cond.L.Unlock()

	ch <- dir
}

func doFindDir(dir string) {
	files, dirs := listDir(dir)

	if *showDirs {
		printer.PrintStrings(dirs)
	}
	if *showFiles {
		printer.PrintStrings(files)
	}
	for _, e := range dirs {
		findDir(e)
	}
}

func listDir(dir string) (files, dirs []string) {
	d, err := os.Open(dir)
	if err != nil {
		common.Warnf("Unable to open %s\n", dir)
		return nil, nil
	}
	defer d.Close()

	children, err := d.Readdirnames(-1)
	if err != nil {
		common.Warnf("Unable to readdir %s\n", dir)
		return nil, nil
	}

	sort.Strings(children)

	files = make([]string, 0, len(children))
	dirs = make([]string, 0, len(children))

	for _, c := range children {
		p := path.Join(dir, c)
		s, err := os.Stat(p)
		if err != nil {
			common.Warnf("Unable to stat %s\n", p)
			continue
		}
		if s.IsDir() {
			dirs = append(dirs, p)
		} else {
			files = append(files, p)
		}
	}
	return
}
