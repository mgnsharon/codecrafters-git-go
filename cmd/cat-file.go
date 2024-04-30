package cmd

import (
	"fmt"
	"os"

	"github.com/fred1268/go-clap/clap"
)

type CatFileParams struct {
	PrettyPrint bool `clap:"--pretty,-p"`
	Hash  []string `clap:"trailing"`
}

func CatFile() {
	var params CatFileParams
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
	obj.PrintContent()
}