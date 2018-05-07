package shescapecommon

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/omakoto/go-common/src/textio"
)

func ShescapeNoNewline(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(shell.Escape(arg))
	}
}

func ShescapeStdin(files []string) {
	err := textio.ReadFiles(files, func(line []byte, lineNo int, filename string) error {
		line, nl := textio.Chomped(line)
		fmt.Print(string(shell.EscapeBytes(line)))
		if nl {
			fmt.Print("\n")
		}
		return nil
	})
	common.Check(err, "Unable to read file")
}

func UnshescapeNoNewline(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(shell.Unescape(arg))
	}
}

func UnshescapeStdin(files []string) {
	err := textio.ReadFiles(files, func(line []byte, lineNo int, filename string) error {
		fmt.Print(string(shell.UnescapeBytes(textio.Chomp(line))), "\n")
		return nil
	})
	common.Check(err, "Unable to read file")
}
