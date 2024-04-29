package cmd

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	"github.com/fred1268/go-clap/clap"
)

type params struct {
	PrettyPrint bool `clap:"--pretty,-p"`
	Hash  []string `clap:"trailing"`
}

func CatFile() {
	var params params
	_, err := clap.Parse(os.Args[2:], &params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %s\n", err)
		os.Exit(1)
	}
	if len(params.Hash) != 1 {
		fmt.Fprintf(os.Stderr, "Invalid Hash\n")
		os.Exit(1)
	}
	
	hash := params.Hash[0]
	fp := fmt.Sprint(".git", string(os.PathSeparator), "objects", string(os.PathSeparator), hash[:2], string(os.PathSeparator), hash[2:])
	
	cf, err := os.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}
	
	uf, err := zlib.NewReader(bytes.NewReader(cf))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decompressing file: %s\n", err)
		os.Exit(1)
	}
	defer uf.Close()
	
	data, err := io.ReadAll(uf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}
	
	filedata := bytes.Split(data, []byte("\x00"))
	/* header := bytes.Split(filedata[0], []byte(" "))
	kind := header[0]
	size := header[1]
	fmt.Println("kind:", string(kind), "size:", string(size)) */
	content := filedata[1]
	io.Copy(os.Stdout, bytes.NewReader(content))

}