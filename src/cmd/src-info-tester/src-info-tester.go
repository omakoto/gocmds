package main

import (
	"fmt"

	"github.com/omakoto/go-common/src/common"
)

func main() {
	file, line := common.GetSourceInfo()
	fmt.Printf("OK: %s, %d\n", file, line)
}
