///bin/true; exec /usr/bin/env go run "$0" "$@"

package main

import (
    "fmt"
    "flag"
    "path/filepath"
)

func main() {
    flag.Parse()

    for _, file := range flag.Args() {
        abs, err := filepath.Abs(file)
        if err != nil {
            panic(err)
        }
        fmt.Print(abs);
        fmt.Print("\n");
    }
}
