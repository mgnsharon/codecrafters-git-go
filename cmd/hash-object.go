package cmd

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/storage"
	"github.com/fred1268/go-clap/clap"
)

type params struct {
	Write bool `clap:"--write,-w"`
	FileName  []string `clap:"trailing"`
}

func HashObject() {
	var params params
	_, err := clap.Parse(os.Args[2:], &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err)
		os.Exit(1)
	}
	if len(params.FileName) != 1 {
		fmt.Fprintf(os.Stderr, "Invalid Filename\n")
		os.Exit(1)
	}
	
	obj := storage.CreateBlob(params.FileName[0])
	fmt.Println(string(obj.ObjectHash))
	if params.Write {
		obj.WriteObject()
	}
	
}