package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/russross/blackfriday/v2"
)

func main() {
	// Read from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	// Convert markdown to HTML
	output := blackfriday.Run(input)

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "mdpreview-*.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary file: %v\n", err)
		os.Exit(1)
	}
	// defer os.Remove(tmpfile.Name()) // Clean up the file afterwards

	// Write the HTML to the temporary file
	if _, err := tmpfile.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to temporary file: %v\n", err)
		os.Exit(1)
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing temporary file: %v\n", err)
		os.Exit(1)
	}

	// Open the file in the browser
	cmd := exec.Command("google-chrome", tmpfile.Name())
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening browser: %v\n", err)
		os.Exit(1)
	}
}
