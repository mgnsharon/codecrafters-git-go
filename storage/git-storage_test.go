package storage

import (
	"bytes"
	"os"
	"path"
	"slices"
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

func TestParseTree(t *testing.T) {
	cdProjectRoot(t)
	testdata := []TreeFile{
		{ "100644", "gitattributes", "176a458f94e0ea5272ce67c36bf30b6be9caf623", TreeContentTypeBlob },
		{ "040000", "test_dir_1", "b31be178b740a3e0fe91468d170000a20a14a269", TreeContentTypeTree },
		{ "040000", "test_dir_2", "8816277598bb0417d1ea4fb40e1a6a487e53b455", TreeContentTypeTree },
		{ "100644", "test_file_1.txt", "3b18e512dba79e4c8300dd08aeb37f8e728b8dad", TreeContentTypeBlob },
	}
	tcs := []struct {
		name string
		hash string
		expected []TreeFile

	}{
		{ "testdata", "abc9ef84782c23419bd9c61f93a1352a26f99ced",  testdata},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			obj := ReadFromHash(tc.hash)
			actual := obj.ParseTree()
			if slices.CompareFunc(actual, tc.expected, func(i, j TreeFile) int {
				if i.Name == j.Name && i.Hash == j.Hash && i.Mode == j.Mode && i.Type == j.Type {
					return 0
				}
				return -1
			}) != 0 {
				t.Errorf("Expected %d, got %d", len(tc.expected), len(actual))
			}
		})
	}
}

func TestCreateTree(t *testing.T) {
	cdProjectRoot(t)
	tcs := []struct {
		name string
		dirname string
		expected string
	}{
		{ "testdata", "storage/testdata", "abc9ef84782c23419bd9c61f93a1352a26f99ced" },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := CreateTree(tc.dirname, []byte{})
			if tc.expected != string(actual.ObjectHash) {
				t.Errorf("Expected %s, got %s", tc.expected, actual.ObjectHash)
			}
		})
	}
}
func TestCreateSha1Hash(t *testing.T) {
	tcs := []struct {
		name string
		input []byte
		expected []byte
	}{
		{ "gitattributes", []byte("176a458f94e0ea5272ce67c36bf30b6be9caf623"), []byte("\x17jE\x8f\x94\xe0\xeaRr\xceg\xc3k\xf3\vk\xe9\xca\xf6#") },
		
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := createSha1Hash(tc.input)
			if !bytes.Equal(tc.expected, actual) {
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