package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	width, height, err := term.GetSize(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "termsize: Failed to get terminal size: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d %d\n", width, height)
}
