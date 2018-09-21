package main

import (
	"os"
	//"path/filepath"
	//"strings"
	"fmt"
	"io"
)

func dirTree(out io.Writer, path string, printFiles bool) (error error) {
	file, err := os.Open(path)

	if err != nil {
		panic("Can not open directory " + path)
	}
	defer file.Close()

	names, _ := file.Readdirnames(0)

	for _, name := range names {
		fmt.Println(name)
	}

	fmt.Println(os.Args)
	fmt.Println(path)
	fmt.Println(printFiles)
	return error
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
