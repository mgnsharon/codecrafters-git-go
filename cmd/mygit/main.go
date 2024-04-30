package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}
	
	switch command := os.Args[1]; command {
	case "init":
		cmd.GitInit()
	case "cat-file":
		cmd.CatFile()
	case "hash-object":
		cmd.HashObject()
	case "ls-tree":
		cmd.LsTree()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
