package cmd

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/storage"
	"github.com/fred1268/go-clap/clap"
)

type LsTreeParams struct {
	NameOnly bool `clap:"--name-only"`
	Hash  []string `clap:"trailing"`
}

func LsTree() {
	var params LsTreeParams
	_, err := clap.Parse(os.Args[2:], &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err)
		os.Exit(1)
	}
	if len(params.Hash) != 1 {
		fmt.Fprintf(os.Stderr, "Invalid Hash\n")
		os.Exit(1)
	}

	files := getFilesFromHash(params.Hash[0])
	
	fmt.Print(getLsTreeOut(files, params.NameOnly))
}

func getFilesFromHash(hash string) []storage.TreeFile {
	obj := storage.ReadFromHash(hash)
	files := obj.ParseTree()

	return files
}

func getLsTreeOut(files []storage.TreeFile, nameOnly bool) string {
	out := ""
	for _, file := range files {
		if nameOnly {
			out += file.Name + "\n"
		} else {
			out += fmt.Sprintf("%s %s %s    %s\n", file.Mode, file.Type, file.Hash, file.Name)
		}
	}
	return out
}
