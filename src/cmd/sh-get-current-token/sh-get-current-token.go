// Intended to be used with "bind -x" on bash.
// Prints the current token in stdout.
// Example: When the command line is "less  /etc/fstab", and when the cursor is at...
// less  /etc/fstab
// ^1 ^2   ^3  ^4
// - ^1:
//      With -f: "less"
//        No -f: ""
// - ^2:
//      With -f: "less"
//        No -f: "less"
// - ^3:
//      With -f: "/etc/fstab"
//        No -f: "/e"
// - ^4:
//      With -f: "/etc/fstab"
//        No -f: "/etc/fstab"

package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/pborman/getopt/v2"
)

var (
	full     = getopt.BoolLong("full", 'f', "Print full token rather than pre-cursor part")
	unescape = getopt.BoolLong("unescape", 'u', "Unescape token")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	sh := shell.MustGetSupportedProxy()

	commandLine, pos := sh.GetCommandLine()
	tokens := sh.Split(commandLine)

	s, ok := findToken(tokens, pos, *full, *unescape)
	if ok {
		fmt.Println(s)
		return 0
	}
	return 1
}

func findToken(tokens []shell.Token, pos int, full, unescae bool) (string, bool) {
	common.Debugf("Pos=%d\n", pos)
	common.Dump("Tokens=", tokens)

	// 0     10    20   30
	// aaa   bbb   ccc  ddd
	//           ^--- 18

	for _, token := range tokens {
		if token.Index > pos {
			break
		}
		if pos > token.Index+len(token.Word) {
			continue
		}
		var t string
		if full {
			t = token.Word
		} else {
			t = token.Word[0 : pos-token.Index]
		}
		if unescae {
			t = shell.Unescape(t)
		}
		return t, true
	}
	return "", false
}
