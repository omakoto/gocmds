///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"flag"
	"fmt"
	"github.com/d4l3k/go-pcre"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/textio"
)

func main() {
	common.RunAndExit(realMain)
}

func realMain() int {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("Usage: %s PATTERN [FILES...]\n", common.MustGetBinName())
		return 1
	}

	pat := flag.Args()[0]

	re, err := pcre.Compile(pat, 0)
	common.Check(err, "Compile() failed")

	found := false
	textio.ReadFiles(flag.Args()[1:], func(line []byte, lineNo int, filename string) error {
		if re.Matcher(line, 0).Matches() {
			fmt.Printf("%s:%d: %s\n", filename, lineNo, textio.Chomp(line))
		}
		return nil
	})
	if found {
		return 0
	}
	return 1
}
