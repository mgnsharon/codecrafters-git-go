package storage

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
)

type ObjectKind string

type StorageWriter interface {
	WriteObject()
}
type TreeContentType string

type TreeFile struct {
	Mode string
	Name string
	Hash string
	Type TreeContentType
}

const (
	ObjectKindBlob   ObjectKind = "blob"
	ObjectKindCommit ObjectKind = "commit"
	ObjectKindTree   ObjectKind = "tree"
	TreeContentTypeBlob TreeContentType = "blob"
	TreeContentTypeTree TreeContentType = "tree"
)

type GitStorage struct {
	Kind ObjectKind
	Size int
	Content []byte
	ObjectHash []byte
}

type GitBlobStorage struct {
	GitStorage
}

func ReadFromHash(hash string) *GitStorage {
	fp := path.Join(".git", "objects", hash[:2], hash[2:])
	
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

	Kind, Size := parseHeader(data)
	Content := data[bytes.IndexByte(data, byte(0))+1:]
	ObjectHash := []byte(hash)
	return &GitStorage{Kind, Size, Content, ObjectHash}
}

func parseHeader(data []byte) (ObjectKind, int) {
	i := bytes.IndexByte(data, byte(0))
	h := bytes.Split(data[:i], []byte(" "))
	var kind ObjectKind
	switch string(h[0]) {
	case "blob":
		kind = ObjectKindBlob
	case "commit":
		kind = ObjectKindCommit
	case "tree":
		kind = ObjectKindTree
	default:
		fmt.Fprintf(os.Stderr, "Unknown object kind: %s\n", kind)
		os.Exit(1)	
	}
	
	size, err := strconv.Atoi(string(h[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting size: %s\n", err)
		os.Exit(1)
	}
	return kind, size
}

func (g *GitStorage) WriteObject() {
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

func CreateBlob(file string) *GitStorage {
	obj := &GitStorage{}
	
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

func (g *GitStorage) computeObjectHash() {
	data := []byte(fmt.Sprintf("%s %d\x00", g.Kind, g.Size))
	data = append(data, g.Content...)
	
	hash := sha1.New()
	hash.Write(data)
	g.ObjectHash = []byte(fmt.Sprintf("%x", hash.Sum(nil)))
	
}

func(g *GitStorage) ParseTree() []TreeFile {
	files := []TreeFile{}
	content := io.Reader(bytes.NewReader(g.Content))	
	var b bytes.Buffer
	io.Copy(&b, content)
	for entry, err := b.ReadBytes(byte(0)); err == nil; {
		w := bytes.Split(entry, []byte(" "))
		Mode := string(w[0])
		if strings.Index(Mode, "4") == 0 {
			Mode = "040000"
		}
		Name := strings.Trim(string(w[1]), "\x00")
		
		var sha [20]byte
		b.Read(sha[:])
		Hash := fmt.Sprintf("%x", sha)
		Type := TreeContentTypeBlob
		if Mode == "040000" {
			Type = TreeContentTypeTree
		}
		files = append(files, TreeFile{Mode, Name, Hash, Type})
		entry, err = b.ReadBytes(byte(0))
	}

	
	slices.SortFunc(files, func(i, j TreeFile) int {
		if i.Name < j.Name {
			return -1
		}
		if i.Name > j.Name {	
			return 1
		}
		return 0
	})
	return files

}

func CreateTree(d string) *GitStorage {
	if d == "" {
		d = "."
	}
	c := []byte{}
	// iterate over the files in the current directory
	 directories, err := os.ReadDir(d)
	 if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %s\n", err)
		os.Exit(1)
	}
	for _, f := range directories {
		if f.Name() == ".git" {
			continue
		}
		
		if f.IsDir() {
			obj := CreateTree(path.Join(d, f.Name()))
			c = append(c, []byte(fmt.Sprintf("040000 %s\x00", f.Name()))...)
			c = append(c, createHash([]byte(obj.ObjectHash))...)
		} else {
			obj := CreateBlob(path.Join(d, f.Name()))
			c = append(c, []byte(fmt.Sprintf("100644 %s\x00", f.Name()))...)
			c = append(c, createHash([]byte(obj.ObjectHash))...)
		}
	}
		// create a Blob for each file
		// and create a Tree for each directory Recursively

		tree := &GitStorage{}
		tree.Kind = ObjectKindTree
		tree.Size = len(c)
		tree.Content = c
		tree.computeObjectHash()
		tree.WriteObject()
		return tree
	}

	func createHash(data []byte) []byte {
		sha := sha1.New()
		sha.Write([]byte(data))
		hash := fmt.Sprintf("%x", sha.Sum(nil))
		dhash, err := hex.DecodeString(hash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding hash: %s\n", err)
			os.Exit(1)
		}
		return dhash
	}