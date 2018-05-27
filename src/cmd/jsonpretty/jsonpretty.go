package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/omakoto/bashcomp"
	"github.com/omakoto/go-common/src/common"
	"io/ioutil"
	"os"
	"strings"
)

var (
	indent = flag.Int("indent", 2, "Number of spaces for indentation")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n  %s: JSON pretty printer\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	bashcomp.HandleBashCompletion()

	bytes, err := ioutil.ReadAll(os.Stdin)
	common.Check(err, "ReadAll failed")

	var f interface{}
	err = json.Unmarshal(bytes, &f)
	common.Check(err, "json.Unmarshal failed")

	// mlib.DebugDump(f)
	// fmt.Println(f)

	pretty, err := json.MarshalIndent(f, "", strings.Repeat(" ", *indent))
	common.Check(err, "json.MarshalIndent failed")

	fmt.Printf("%s\n", pretty)
}
