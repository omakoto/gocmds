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
	"strings"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	sh := shell.GetSupportedProxy()
	common.OrFatalf(sh != nil, "Unsupported shell.\n")

	commandLine, pos := sh.GetCommandLine()
	tokens := sh.Split(commandLine)

	var command shell.Token
	first := true
	for _, t := range tokens {
		if t.Index > pos {
			break
		}
		if strings.ContainsAny(string(t.Word[0]), "|&;{}()!") {
			first = true
			continue
		}
		if first {
			command = t
			first = false
		}
	}
	if command.Word != "" {
		fmt.Print(command.Word, "\n")
	}
	return 0
}
