// multiplex is a command that reads multiple files concurrently and prints their contents line by line to stdout.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

func main() {
	files := os.Args[1:]
	var wg sync.WaitGroup
	var mu sync.Mutex

	var numSuccess atomic.Int32

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: multiplex <file1> <file2> ...")
		os.Exit(3)
	}

	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", file, err)
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				mu.Lock()
				fmt.Print(scanner.Text())
				fmt.Print("\n")
				mu.Unlock()
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", file, err)
				return
			}
			numSuccess.Add(1)
		}(file)
	}
	wg.Wait()

	var ns = int(numSuccess.Load())
	status := 0
	if ns == 0 {
		status = 2
	} else if ns < len(files) {
		status = 1
	}
	os.Exit(status)
}
