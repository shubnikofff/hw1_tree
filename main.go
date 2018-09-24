package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func readDir(path string) (err error, files []os.FileInfo) {
	file, err := os.Open(path)
	files, err = file.Readdir(0)
	file.Close()
	return err, files
}

func printDir(out io.Writer, nodes []os.FileInfo, path string, prefix []string) (err error) {
	fmt.Fprintf(out, "%s", strings.Join(prefix, ""))

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name() < nodes[j].Name()
	})

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node.Name())
		if node.IsDir() {
			nextDir := filepath.Join(path, node.Name())
			_, files := readDir(nextDir)
			printDir(out, files, filepath.Join(path, node.Name()), append(prefix, "\t"))
		}

		return nil
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node.Name())
	if node.IsDir() {
		nextDir := filepath.Join(path, node.Name())
		_, files := readDir(nextDir)
		printDir(out, files, filepath.Join(path, node.Name()), append(prefix, "|\t"))
	}

	printDir(out, nodes[1:], path, prefix)

	return err
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	_, files := readDir(path)
	printDir(out, files, path, []string{})
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
