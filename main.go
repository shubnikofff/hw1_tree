package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	name string
}

type Directory struct {
	Node
	children []fmt.Stringer
}

type File struct {
	Node
	size int64
}

func (file File) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	}
	return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
}

func (directory Directory) String() string {
	return directory.name
}

func readDir(path string, nodes []fmt.Stringer) []fmt.Stringer {
	file, _ := os.Open(path)
	files, _ := file.Readdir(0)
	file.Close()

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, info := range files {
		var newNode fmt.Stringer

		if info.IsDir() {
			newNode = Directory{
				Node:     Node{info.Name()},
				children: readDir(filepath.Join(path, info.Name()), []fmt.Stringer{}),
			}
		} else {
			newNode = File{Node{info.Name()}, info.Size()}
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}

func printDir(out io.Writer, nodes []fmt.Stringer, prefix []string) {
	fmt.Fprintf(out, "%s", strings.Join(prefix, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if directory, ok := node.(Directory); ok {
			printDir(out, directory.children, append(prefix, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if directory, ok := node.(Directory); ok {
		printDir(out, directory.children, append(prefix, "│\t"))
	}

	printDir(out, nodes[1:], prefix)
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	nodes := readDir(path, []fmt.Stringer{})
	printDir(out, nodes, []string{})

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
