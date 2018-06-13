package main

import (
	"bytes"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/ungerik/go-dry"
	"os"
	"strings"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	b := &bytes.Buffer{}

	args := os.Args[1:]

	f := make([]string, 0)
	a := make([]string, 0)

	i := 0
	for ; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") || args[i] == "--" {
			break
		}
		f = append(f, args[i])
	}

	for ; i < len(args); i++ {
		a = append(a, args[i])
	}
	b.WriteString("all_flags=(")
	b.WriteString(strings.Join(dry.StringMap(func(s string) string {
		return shell.Escape(s)
	}, f), " "))
	b.WriteString(")\n")

	b.WriteString("set -- ")
	b.WriteString(strings.Join(dry.StringMap(func(s string) string {
		return shell.Escape(s)
	}, a), " "))
	b.WriteString("\n")
	os.Stdout.Write(b.Bytes())
	return 0
}
