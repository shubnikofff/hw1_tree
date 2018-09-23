package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	file, err := os.Open(path)
	dir, err := file.Readdir(0)
	file.Close()

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})

	var b strings.Builder
	for range strings.Split(path, "/") {
		b.WriteString("|\t")
	}

	for _, info := range dir {
		if info.IsDir() {
			fmt.Fprintf(out, "%s%s\n", b.String(), info.Name())
			dirTree(out, filepath.Join(path, info.Name()), printFiles)
		} else {
			if printFiles {
				fmt.Fprintf(out, "%s%s (%db)\n",b.String(), info.Name(), info.Size())
			}
		}
	}

	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
