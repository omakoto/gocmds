// Intended to be used with "bind -x" on bash.
// Replace the current token in the command line with a given string.
package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/omakoto/go-common/src/textio"
	"github.com/pborman/getopt/v2"
	"io"
)

var (
	index = getopt.BoolLong("index", 'i', "Insert rather than replace")
	null  = getopt.BoolLong("null", '0', "Escape given word before replace")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	eol := byte('\n')
	if *null {
		eol = 0
	}
	printWords(textio.BufferedStdout, eol, *index)

	return 0
}

func printWords(out io.Writer, eol byte, showIndex bool) {
	sh := shell.MustGetSupportedProxy()
	commandLine, _ := sh.GetCommandLine()
	tokens := sh.Split(commandLine)

	eolBytes := append(make([]byte, 0), eol)

	for _, t := range tokens {
		if showIndex {
			fmt.Fprintf(out, "%d ", t.Index)
		}
		out.Write([]byte(t.Word))
		out.Write(eolBytes)
	}
}
