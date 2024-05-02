package cmd

import (
	"fmt"

	"github.com/codecrafters-io/git-starter-go/storage"
)

func WriteTree() {
	obj := storage.CreateTree(".")
	fmt.Println(obj.ObjectHash)
}