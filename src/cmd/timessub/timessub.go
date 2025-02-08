package main

/*
timessub calculates user/system differences obtained from bash's `times` command.

Usage:
times >/dev/shm/$$-pre.txt &&
some slow command &&
times >/dev/shm/$$-post.txt &&
timessub $(cat /dev/shm/$$-post.txt) $(cat /dev/shm/$$-pre.txt)

or,
timessub -f /dev/shm/$$-post.txt /dev/shm/$$-pre.txt

*/

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/must"
	"github.com/pborman/getopt"
)

/* Sample times output
$ times
0m0.235s 0m0.502s # shell:    user, sys
0m5.911s 0m7.013s # children: user, sys

    Prints the accumulated user and system times for the shell and all of its
    child processes.

*/

var (
	fromFile = getopt.BoolLong("file", 'f', "read old / new values from files")
)

func main() {
	common.RunAndExit(realMain)
}

func parseTime(time string) float64 {
	if !strings.HasSuffix(time, "s") {
		common.Fatalf("Invalid input '%s' -- no 's' suffix", time)
	}
	time = time[0 : len(time)-1]
	fields := strings.Split(time, "m")

	if len(fields) != 2 {
		common.Fatalf("Invalid input '%s' -- no 'm'", time)
	}

	seconds := must.Must2(strconv.ParseFloat(fields[1], 64))
	minutes := must.Must2(strconv.Atoi(fields[0]))

	return float64(minutes*60) + seconds
}

func realMain() int {
	getopt.Parse()

	args := getopt.Args()
	if *fromFile {
		return doFiles(args)
	} else {
		return doArgs(args)
	}
}

func print(newUser, newSys, oldUser, oldSys float64) {
	fmt.Printf("times_user=%0.6f\ntimes_sys=%0.6f\n", (newUser - oldUser), (newSys - oldSys))
}

func doArgs(args []string) int {
	if len(args) != 8 {
		fmt.Fprintf(os.Stderr, "Usage: timessub $(times [new]) $(times [old])\n")
		return 1
	}
	newUser := parseTime(args[2])
	newSys := parseTime(args[3])

	oldUser := parseTime(args[6])
	oldSys := parseTime(args[7])

	print(newUser, newSys, oldUser, oldSys)

	return 0
}

func readFile(file string) (float64, float64) {
	f := must.Must2(os.Open(file))
	defer f.Close()

	s := string(must.Must2(io.ReadAll(f)))

	fields := strings.Fields(s)
	if len(fields) < 4 {
		common.Fatalf("Unable to parse '%s': content='%s'", file, s)
	}
	return parseTime(fields[2]), parseTime(fields[3])
}

func doFiles(args []string) int {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: timessub -f [new-times-file] [old-times-file]\n")
		return 1
	}
	newUser, newSys := readFile(args[0])
	oldUser, oldSys := readFile(args[1])

	print(newUser, newSys, oldUser, oldSys)

	return 0
}
