package shescapecommon

import (
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/shell"
	"github.com/omakoto/go-common/src/textio"
)

func ShescapeNoNewline(args []string) {
	for i, arg := range args {
		if i > 0 {
			textio.BufferedStdout.WriteByte(' ')
		}
		textio.BufferedStdout.WriteString(shell.Escape(arg))
	}
}

func ShescapeStdin(files []string) {
	err := textio.ReadFiles(files, func(line []byte, lineNo int, filename string) error {
		line, nl := textio.Chomped(line)
		textio.BufferedStdout.Write(shell.EscapeBytes(line))
		textio.BufferedStdout.Write(nl)
		return nil
	})
	common.Check(err, "Unable to read file")
}

func UnshescapeNoNewline(args []string) {
	for i, arg := range args {
		if i > 0 {
			textio.BufferedStdout.WriteByte(' ')
		}
		textio.BufferedStdout.WriteString(shell.Unescape(arg))
	}
}

func UnshescapeStdin(files []string) {
	err := textio.ReadFiles(files, func(line []byte, lineNo int, filename string) error {
		line, nl := textio.Chomped(line)
		textio.BufferedStdout.Write(shell.UnescapeBytes(line))
		textio.BufferedStdout.Write(nl)
		return nil
	})
	common.Check(err, "Unable to read file")
}
