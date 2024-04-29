package main

type ObjectHash string

type ObjectKind string

const (
	ObjectKindBlob   ObjectKind = "blob"
	ObjectKindCommit ObjectKind = "commit"
	ObjectKindTree   ObjectKind = "tree"
)

type GitObject struct {
	Kind ObjectKind
	ObjectHash	ObjectHash
}