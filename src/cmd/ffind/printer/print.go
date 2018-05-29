package printer

import (
	"bufio"
	"os"
	"sync"
)

var (
	out = bufio.NewWriterSize(os.Stdout, 64*1024)
	sem = sync.Mutex{}
)

func Flush() {
	out.Flush()
}

func PrintStrings(strs []string) {
	sem.Lock()
	defer sem.Unlock()

	for _, s := range strs {
		out.WriteString(s)
		out.WriteByte('\n')
	}
}
