// Intended to be used with "bind -x" on bash.
// Prints the executable command name from the command line.
// If the command line is a pipeline, then print the one that's right before the cursor.
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
	sh := shell.MustGetSupportedProxy()

	commandLine, pos := sh.GetCommandLine()
	tokens := sh.Split(commandLine)
	commandName := findCommandName(tokens, pos)

	if commandName != "" {
		fmt.Print(commandName, "\n")
	}
	return 0
}

func findCommandName(tokens []shell.Token, pos int) string {
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
	return command.Word
}
