package main

import (

	"github.com/d4l3k/go-pcre"
	"fmt"
)

func main() {
	re := pcre.MustCompile("test", 0)
	m := re.Matcher([]byte("abctest"), 0)
	fmt.Printf("%s", m.Matches())
}