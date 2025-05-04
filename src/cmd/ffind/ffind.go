package main

import (
	"os"
	"path"
	"runtime"
	"sort"
	"sync"

	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/ffind/printer"
	"github.com/pborman/getopt/v2"
)

const maxPara = 4

func defaultPara() int {
	v := runtime.NumCPU() * 4
	if v < maxPara {
		return v
	}
	return maxPara
}

var (
	showFiles      = getopt.BoolLong("file", 'f', "Print files only")
	showDirs       = getopt.BoolLong("dir", 'd', "Print directories only")
	followSymlinks = getopt.BoolLong("symlink", 'L', "Follow symlinks")

	quiet = getopt.BoolLong("quiet", 'q', "Don't show warnings")

	para = getopt.IntLong("para", 'j', defaultPara(), "Number of goroutines")

	ch = make(chan string, 100*1024)

	numBacklog = 0
	cond       = sync.NewCond(&sync.Mutex{})

	statFunc func(string) (os.FileInfo, error)
)

func main() {
	common.RunAndExit(realMain)
}

func startWorker() {
	go func() {
		files := make([]string, 0, 1024)
		dirs := make([]string, 0, 1024)

		for {
			dir := <-ch
			//common.Debugf("Pop:  %s\n", dir)
			doFindDir(dir, files, dirs)

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

func realMain() int {
	getopt.Parse()

	common.Quiet = *quiet

	if !(*showFiles || *showDirs) {
		*showFiles = true
		*showDirs = true
	}

	if *followSymlinks {
		statFunc = os.Stat
	} else {
		statFunc = os.Lstat
	}

	common.Debugf("-j=%d\n", *para)

	for i := 0; i < *para; i++ {
		startWorker()
	}

	if len(getopt.Args()) == 0 {
		dir, err := os.Getwd()
		common.Check(err, "Getwd failed")

		pushDir(dir)
	} else {
		for _, dir := range getopt.Args() {
			pushDir(dir)
		}
	}

	cond.L.Lock()
	if numBacklog > 0 {
		cond.Wait()
	}
	cond.L.Unlock()

	return 0
}

func pushDir(dir string) {
	cond.L.Lock()
	//common.Debugf("Push: %s [%d]\n", dir, numBacklog)
	numBacklog++
	cond.L.Unlock()

	for {
		select {
		case ch <- dir:
			return
		default:
			startWorker() // Just create more workers when buffer is full.
		}
	}
}

func clearStringSlice(s []string) []string {
	for i := 0; i < len(s); i++ {
		s[i] = ""
	}
	return s[:0]
}

func doFindDir(dir string, files, dirs []string) {
	files = clearStringSlice(files)
	dirs = clearStringSlice(dirs)

	files, dirs = listDir(dir, files, dirs)

	if *showDirs {
		printer.PrintStrings(dirs)
	}
	if *showFiles {
		printer.PrintStrings(files)
	}
	for _, e := range dirs {
		pushDir(e)
	}
}

func listDir(dir string, files, dirs []string) ([]string, []string) {
	d, err := os.Open(dir)
	if err != nil {
		common.Warnf("Unable to open %s\n", dir)
		return files, dirs
	}
	defer d.Close()

	children, err := d.Readdirnames(-1)
	if err != nil {
		common.Warnf("Unable to readdir %s\n", dir)
		return files, dirs
	}

	sort.Strings(children)

	for _, c := range children {
		p := path.Join(dir, c)
		common.Debugf("  %s\n", p)
		s, err := statFunc(p)
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
	return files, dirs
}
