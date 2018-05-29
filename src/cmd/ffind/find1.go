package main
//
//import (
//	"bufio"
//	"bytes"
//	"github.com/omakoto/go-common/src/common"
//	"github.com/pborman/getopt/v2"
//	"os"
//	"path"
//	"sort"
//	"strconv"
//)
//func getDefaultCachedDir() string {
//  return path.Join(os.Getenv("HOME"), ".ffind")
//}
//
//var (
//	cacheDir = getopt.StringLong("cache-dir", 'c', getDefaultCachedDir(), "Specify cache directory")
//)
//
//func listDir(dir string) (files, dirs []string, err error) {
//	cacheFile := path.Join(*cacheDir, dir[1:], "ffind-cache.txt")
//
//	// If already cached, just return it.
//	avail, files, dirs := listCachedDir(cacheFile, dir)
//	if avail {
//		return
//	}
//
//	// Not cached.
//	d, err := os.Open(dir)
//	if err != nil {
//		common.Warnf("Unable to open %s\n", dir)
//		return nil, nil, nil
//	}
//	defer d.Close()
//
//	children, err := d.Readdirnames(-1)
//	if err != nil {
//		common.Warnf("Unable to readdir %s\n", dir)
//		return nil, nil, nil
//	}
//
//	sort.Strings(children)
//
//	files = make([]string, 0, len(children))
//	dirs = make([]string, 0, len(children))
//
//	for _, c := range children {
//		p := path.Join(dir, c)
//		s, err := os.Stat(p)
//		if err != nil {
//			common.Warnf("Unable to stat %s\n", p)
//			continue
//		}
//		if s.IsDir() {
//			dirs = append(dirs, p)
//		} else {
//			files = append(files, p)
//		}
//	}
//
//	mustWriteCache(cacheFile, files, dirs)
//
//	return
//}
//
//func listCachedDir(cacheFile, dir string) (avail bool, files, dirs []string) {
//	// See if valid cache.
//	scache, err := os.Stat(cacheFile)
//	if err != nil {
//		common.Debugf("Cache %s doesn't exist", cacheFile)
//		return false, nil, nil
//	}
//	sdir, err := os.Stat(dir)
//	if err != nil {
//		common.Debugf("Dir %s doesn't exist\n", dir)
//		return false, nil, nil
//	}
//	if scache.ModTime().Before(sdir.ModTime()) {
//		common.Debugf("Cache %s older than %s\n", cacheFile, dir)
//		return false, nil, nil
//	}
//	in, err := os.Open(cacheFile)
//	common.Check(err, "unable to open cache file\n")
//	defer in.Close()
//
//	// Cache valid, read it and return.
//
//	const errFormat = "cache file format error"
//
//	b := bufio.NewReader(in)
//
//	readInt := func() int {
//		n, err := b.ReadBytes('\n')
//		common.Check(err, errFormat)
//
//		i, err := strconv.ParseInt(string(bytes.TrimRight(n, "\n")), 10, 32)
//		common.Check(err, errFormat)
//
//		return int(i)
//	}
//
//	numFiles := readInt()
//	numDirs := readInt()
//
//	files = make([]string, 0, numFiles)
//	dirs = make([]string, 0, numDirs)
//
//	for i := 0; i < int(numFiles); i++ {
//		n, err := b.ReadBytes('\n')
//		common.Check(err, errFormat)
//
//		files = append(files, string(bytes.TrimRight(n, "\n")))
//	}
//
//	for i := 0; i < int(numDirs); i++ {
//		n, err := b.ReadBytes('\n')
//		common.Check(err, errFormat)
//
//		dirs = append(dirs, string(bytes.TrimRight(n, "\n")))
//	}
//
//	return true, files, dirs
//}
//
//func mustWriteCache(cacheFile string, files, dirs []string) {
//	dir := path.Dir(cacheFile)
//	err := os.MkdirAll(dir, 0700)
//	common.Checkf(err, "unable to create cache directory %s", dir)
//
//	out, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0500)
//	common.Checkf(err, "unable to open cache file %s", cacheFile)
//	defer out.Close()
//
//	b := bufio.NewWriter(out)
//
//	b.WriteString(strconv.Itoa(len(files)))
//	b.WriteByte('\n')
//	b.WriteString(strconv.Itoa(len(dirs)))
//	b.WriteByte('\n')
//
//	for _, f := range files {
//		b.WriteString(f)
//		b.WriteByte('\n')
//	}
//
//	for _, d := range dirs {
//		b.WriteString(d)
//		b.WriteByte('\n')
//	}
//	b.Flush()
//
//	common.Debugf("Wrote cache file %s\n", cacheFile)
//}
