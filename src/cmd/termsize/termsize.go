package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

func main() {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "termsize: Failed to get terminal size: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d %d\n", width, height)
}
