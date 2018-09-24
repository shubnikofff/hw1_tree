package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Node struct {
	name  string
	nodes []Node
}

func (node Node) hasNodes() bool {
	return len(node.nodes) > 0
}

func readDir(path string, nodes []Node) []Node {
	file, _ := os.Open(path)
	files, _ := file.Readdir(0)
	file.Close()

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, info := range files {
		var newNode = Node{
			name: info.Name(),
		}

		if info.IsDir() {
			newNode.nodes = readDir(filepath.Join(path, info.Name()), []Node{})
		}

		nodes = append(nodes, newNode)
	}

	return nodes
}

func printDir(out io.Writer, nodes []Node, prefix []string) {
	fmt.Fprintf(out, "%s", strings.Join(prefix, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node.name)
		if node.hasNodes() {
			printDir(out, node.nodes, append(prefix, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node.name)
	if node.hasNodes() {
		printDir(out, node.nodes, append(prefix, "|\t"))
	}

	printDir(out, nodes[1:], prefix)
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	nodes := readDir(path, []Node{})
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
