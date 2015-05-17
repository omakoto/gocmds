///bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
    "fmt"
    "flag"
    "os"
    "path/filepath"
)

func main() {
    flag.Parse()

    cwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    for _, file := range flag.Args() {
        rel, err := filepath.Rel(cwd, file)
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
            continue
        }
        fmt.Print(rel);
        fmt.Print("\n");
    }
}
