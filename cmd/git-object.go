package cmd

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

type ObjectHash string

type ObjectKind string

const (
	ObjectKindBlob   ObjectKind = "blob"
	ObjectKindCommit ObjectKind = "commit"
	ObjectKindTree   ObjectKind = "tree"
)

type GitObject struct {
	Kind ObjectKind
	Size int
	Content []byte
	ObjectHash ObjectHash
}

func (g *GitObject) GetObjectHash() ObjectHash {
	return g.ObjectHash
}

func (g *GitObject) GetKind() ObjectKind {
	return g.Kind
}

func (g *GitObject) SetObjectHash(hash ObjectHash) {
	g.ObjectHash = hash
}

func (g *GitObject) SetKind(kind ObjectKind) {
	g.Kind = kind
}

func ReadFromHash(hash string) *GitObject {
	obj := &GitObject{}
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
	header := bytes.Split(filedata[0], []byte(" "))
	kind := header[0]
	switch string(kind) {
	case "blob":
		obj.Kind = ObjectKindBlob
	case "commit":
		obj.Kind = ObjectKindCommit
	case "tree":
		obj.Kind = ObjectKindTree
	default:
		fmt.Fprintf(os.Stderr, "Unknown object kind: %s\n", kind)
		os.Exit(1)	
	}
	obj.Size, err = strconv.Atoi(string(header[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting size: %s\n", err)
		os.Exit(1)
	}
	
	obj.Content = filedata[1]
	
	return obj
}

func CreateBlob(file string) *GitObject {
	obj := &GitObject{}
	
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	obj.Kind = ObjectKindBlob
	obj.Size = len(data)
	obj.Content = data
	obj.computeObjectHash()
	return obj
}

func (g *GitObject) computeObjectHash() {
	data := []byte(fmt.Sprintf("%s %d\x00", g.Kind, g.Size))
	data = append(data, g.Content...)
	
	hash := sha1.New()
	hash.Write(data)
	g.ObjectHash = ObjectHash(fmt.Sprintf("%x", hash.Sum(nil)))
	
}

func (g *GitObject) WriteObject() {
	data := []byte(fmt.Sprintf("%s %d\x00", g.Kind, g.Size))
	data = append(data, g.Content...)
	
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	
	_, err := w.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}
	w.Close()

	fp := path.Join(".git", "objects", string(g.ObjectHash)[:2], string(g.ObjectHash)[2:])
	
	if err := os.MkdirAll(path.Dir(fp), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(fp, b.Bytes(), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		os.Exit(1)
	}
}

func (g *GitObject) PrintContent() {
	io.Copy(os.Stdout, bytes.NewReader(g.Content))
}

