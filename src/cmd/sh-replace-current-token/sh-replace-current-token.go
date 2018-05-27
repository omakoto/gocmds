// Intended to be used with "bind -x" on bash.
// Replace the current token in the command line with a given string.
package main

import (
	"bytes"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/pborman/getopt/v2"
	"strings"
)

var (
	insert = getopt.BoolLong("insert", 'i', "Insert rather than replace")
	escape = getopt.BoolLong("escape", 'e', "Escape given word before replace")
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	getopt.Parse()

	sh := shell.GetSupportedProxy()
	common.OrFatalf(sh != nil, "Unsupported shell.\n")

	commandLine, pos := sh.GetCommandLine()

	newWord := strings.Join(getopt.Args(), " ")

	commandLine, pos = doTransform(commandLine, pos, newWord, *insert, *escape)

	sh.PrintUpdateCommandLineEvalStr(commandLine, pos)

	return 0
}

func doTransform(original string, pos int, newWord string, insert, escape bool) (string, int) {
	if escape {
		newWord = shell.Escape(newWord)
	}

	if insert {
		original, pos = insertWord(original, pos, newWord)
	} else {
		original, pos = replaceWord(original, pos, newWord)
	}

	return original, pos
}

func insertWord(original string, pos int, newWord string) (string, int) {
	var ret string
	if pos <= len(original) {
		ret = original[0:pos] + newWord + original[pos:]
		pos += len(newWord)
	} else {
		ret = original + " " + newWord
		pos = len(ret)
	}

	return ret, pos
}

func replaceWord(original string, pos int, newWord string) (string, int) {
	tokens := shell.SplitToTokens(original)

	// 0     10    20   30
	// aaa   bbb   ccc  ddd
	//           ^--- 18

	// Find the token at the right position and replace it.
	found := false
	newPos := 0
	for i, token := range tokens {
		if token.Index > pos {
			break
		}
		if pos > token.Index+len(token.Word) {
			continue
		}

		// Adjust the indexes of the following tokens
		lenDelta := len(newWord) - len(token.Word)
		for j := i + 1; j < len(tokens); j++ {
			tokens[j].Index += lenDelta
		}
		tokens[i] = shell.Token{Word: newWord, Index: token.Index}
		newPos = token.Index + len(newWord)

		found = true
		break
	}
	// If not found, fall back to insert.
	if !found {
		return insertWord(original, pos, newWord)
	}

	ret := bytes.Buffer{}
	for _, token := range tokens {
		for ret.Len() < token.Index {
			ret.WriteByte(' ')
		}
		ret.WriteString(token.Word)
	}
	return ret.String(), newPos
}
