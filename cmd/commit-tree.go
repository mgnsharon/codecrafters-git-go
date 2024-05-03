package cmd

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/storage"
	"github.com/fred1268/go-clap/clap"
)

type commitTreeParams struct {
	Parent string `clap:",-p"`
	Message string `clap:",-m"`
	Hash string
}

func CommitTree() {
	var params commitTreeParams
	params.Hash = os.Args[2]
	_, err := clap.Parse(os.Args[3:], &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err)
		os.Exit(1)
	}
	
	commit := storage.CreateCommit(params.Hash, params.Message, params.Parent)
	commit.WriteObject()
	fmt.Println(string(commit.ObjectHash))
	
}


