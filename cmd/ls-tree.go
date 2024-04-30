package cmd

import (
	"fmt"
	"os"

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

	obj := ReadFromHash(params.Hash[0])
	files := obj.ParseTree()
	
	for _, file := range files {
		if params.NameOnly {
			fmt.Println(file.Name)
		} else {
			fmt.Printf("%s %s %s %s\n", file.Mode, file.Type, file.Hash, file.Name)
		}
	}
}