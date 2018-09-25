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

type Node interface {
	fmt.Stringer
}

type Directory struct {
	name     string
	children []Node
}

type File struct {
	name string
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

func readDir(path string, nodes []Node, withFiles bool) (error, []Node) {
	file, err := os.Open(path)
	files, err := file.Readdir(0)
	file.Close()

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, info := range files {
		if !(info.IsDir() || withFiles) {
			continue
		}

		var newNode Node
		if info.IsDir() {
			_, children := readDir(filepath.Join(path, info.Name()), []Node{}, withFiles)
			newNode = Directory{info.Name(), children}
		} else {
			newNode = File{info.Name(), info.Size()}
		}

		nodes = append(nodes, newNode)
	}

	return err, nodes
}

func printDir(out io.Writer, nodes []Node, prefix []string) {
	if len(nodes) == 0 {
		return
	}

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	err, nodes := readDir(path, []Node{}, printFiles)
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
