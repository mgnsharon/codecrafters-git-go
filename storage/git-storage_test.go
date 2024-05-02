package storage

import (
	"bytes"
	"os"
	"path"
	"testing"
)

func TestReadFromHash(t *testing.T) {
	cdProjectRoot(t)
	
	tcs := []struct {
		name string
		hash string
		expected []byte
	}{
		{ "test_file_1.txt", "3b18e512dba79e4c8300dd08aeb37f8e728b8dad", []byte("hello world") },
		{ "gitattributes", "176a458f94e0ea5272ce67c36bf30b6be9caf623", []byte("* text=auto") },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := ReadFromHash(tc.hash)
			if bytes.Equal(tc.expected, actual.Content) {
				t.Errorf("Expected %s, got %s", string(tc.expected), string(actual.Content))
			}
		})
	}
}

func TestCreateBlob(t *testing.T) {
	cdProjectRoot(t)
	tcs := []struct {
		name string
		fileName string
		expected string
	}{
		{ "gitattributes", "storage/testdata/gitattributes", "176a458f94e0ea5272ce67c36bf30b6be9caf623" },
		{ "test_file_1.txt", "storage/testdata/test_file_1.txt", "3b18e512dba79e4c8300dd08aeb37f8e728b8dad" },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateBlob(tc.fileName)
			if tc.expected != string(actual.ObjectHash) {
				t.Errorf("Expected %s, got %s", tc.expected, actual.ObjectHash)
			}
		})
	}
}

func TestCreateTree(t *testing.T) {
	tcs := []struct {
		name string
		dirname string
		expected string
	}{
		{ "testdata", "testdata", "abc9ef84782c23419bd9c61f93a1352a26f99ced" },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateTree(tc.dirname)
			if tc.expected != string(actual.ObjectHash) {
				t.Errorf("Expected %s, got %s", tc.expected, actual.ObjectHash)
			}
		})
	}
}
func TestCreateHash(t *testing.T) {
	tcs := []struct {
		name string
		input []byte
		expected []byte
	}{
		{ "empty string", []byte(""), []byte("e69de29bb2d1d6434b8b29ae775ad8c2e48c5391") },
		{ "hello world", []byte("hello world"), []byte("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c") },
		{ "hello world with newline", []byte("hello world\n"), []byte("6adfb183a4a2c94a2f92dab5ade762a47889a5a1") },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := createHash(tc.input)
			if bytes.Equal(tc.expected, actual) {
				t.Errorf("Expected %s, got %s", tc.expected, actual)
			}
		})
	}
}

func cdProjectRoot(t *testing.T) {
	t.Helper()
	d, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory: %s", err)
	}
	findProjectRoot(t, d)
}

func findProjectRoot(t *testing.T, dir string) {
	t.Logf("Looking for .git in %s", dir)
	fp := path.Join(dir, ".git")

	if _, err := os.Stat(fp); err == nil {
		t.Logf("Found .git in %s", dir)
		os.Chdir(dir)
		return
	}
	findProjectRoot(t, dir + "/..")
}