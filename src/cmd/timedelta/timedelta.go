package main

// timedelta keeps track of the time every time it's executed,
// and prints the delta time since the last execution.
//
// Must pass a unique string as the first argument to distinguish
// multiple "sessions".
//
// Useful for profiling a shell script, with:
// set -x ; PS4='+ $(timedelta '$$') '

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/omakoto/go-common/src/common"
)

func main() {
	common.RunAndExit(realMain)
}

func read(fname string) (time.Time, error) {
	f, err := os.Open(fname)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return time.Time{}, err
	}
	v, err := strconv.Atoi(strings.TrimRight(string(b), "\n"))
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMicro(int64(v)), nil
}

func write(fname string, t time.Time) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0600)
	common.Checke(err)
	defer f.Close()
	fmt.Fprint(f, t.UnixMicro())
}

func realMain() int {

	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: timedelta UNIQUE-KEY (typically $$)\n")
		return 1
	}
	key := args[0]

	now := time.Now()

	file := "/dev/shm/timedelta-" + key + ".tmp"

	last, err := read(file)
	if err != nil {
		// first call
		fmt.Printf("0\n")
	} else {
		fmt.Printf("%.6f\n", now.Sub(last).Seconds())
	}

	write(file, now)

	return 0
}
