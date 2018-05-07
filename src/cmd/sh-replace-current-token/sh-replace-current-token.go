package main

import (
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
		panic("TODO Implement it")
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
