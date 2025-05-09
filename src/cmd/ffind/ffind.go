package main

import (
	"os"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/gocmds/src/cmd/ffind/printer"
	"github.com/pborman/getopt/v2"
)

const maxPara = 16

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
	quiet          = getopt.BoolLong("quiet", 'q', "Don't show warnings")
	para           = getopt.IntLong("para", 'j', defaultPara(), "Number of goroutines")
	noSkip         = getopt.BoolLong("no-skip", 'a', "Don't skip .git, etc")

	excludeDirs []string
	excludeMap  map[string]bool

	skipPatterns []string
	skipRe       *regexp.Regexp

	ch = make(chan string, 100*1024)

	numBacklog = 0
	cond       = sync.NewCond(&sync.Mutex{})

	statFunc    func(string) (os.FileInfo, error)
	includeTest func(fullPath string) bool
	shouldSkip  func(dirName string) bool
)

func init() {
	getopt.FlagLong(&excludeDirs, "exclude", 'x', "directories to exclude in full path")
	getopt.FlagLong(&skipPatterns, "ignore", 'i', "names of the directories to skip")
}

func main() {
	common.RunAndExit(realMain)
}

func mapToTest(m map[string]bool, def bool) func(string) bool {
	if len(m) == 0 {
		return func(string) bool {
			return def
		}
	} else {
		return func(s string) bool {
			return m[s]
		}
	}
}

func initialize() {
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

	shouldSkipForGit := func(dirName string) bool {
		if dirName == ".git" || dirName == ".repo" {
			return true
		}
		return false
	}
	if *noSkip {
		shouldSkipForGit = func(dirName string) bool {
			return false
		}
	}

	shouldSkipForIgnore := func(dirName string) bool {
		return false
	}
	if len(skipPatterns) > 0 {
		skipRe = regexp.MustCompile("(" + strings.Join(skipPatterns, "|") + ")")
		shouldSkipForIgnore = func(dirName string) bool {
			return skipRe.MatchString(dirName)
		}
	}

	shouldSkip = func(dirName string) bool {
		if shouldSkipForGit(dirName) {
			return true
		}
		if shouldSkipForIgnore(dirName) {
			return true
		}

		return false
	}

	excludeMap = make(map[string]bool)

	for _, d := range excludeDirs {
		excludeMap[d] = true
	}

	et := mapToTest(excludeMap, false)

	includeTest = func(fullPath string) bool {
		return !et(fullPath)
	}

	//common.Debugf("-j=%d\n", *para)

}

func realMain() int {
	initialize()

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

	printer.PrintStrings(dirs)
	printer.PrintStrings(files)

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
		common.Debugf("  %s - %s\n", dir, c)
		s, err := statFunc(p)
		if err != nil {
			common.Warnf("Unable to stat %s\n", p)
			continue
		}
		if shouldSkip(c) {
			continue
		}
		if s.IsDir() {
			if *showDirs && includeTest(p) {
				dirs = append(dirs, p)
			}
		} else {
			if *showFiles {
				files = append(files, p)
			}
		}
	}
	return files, dirs
}
