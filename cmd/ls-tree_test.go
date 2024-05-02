package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/codecrafters-io/git-starter-go/storage"
)
const (
	LSTreeOutNameOnly = "gitattributes\ntest_dir_1\ntest_dir_2\ntest_file_1.txt\n"
	LSTreeOutFull = `100644 blob 176a458f94e0ea5272ce67c36bf30b6be9caf623    gitattributes
040000 tree b31be178b740a3e0fe91468d170000a20a14a269    test_dir_1
040000 tree 8816277598bb0417d1ea4fb40e1a6a487e53b455    test_dir_2
100644 blob 3b18e512dba79e4c8300dd08aeb37f8e728b8dad    test_file_1.txt
`
)

func TestGetFilesFromHash(t *testing.T) {
	cdProjectRoot(t)
	tcs := []struct {
		name string
		hash string
		expected int
	}{
		{ "testdata", "abc9ef84782c23419bd9c61f93a1352a26f99ced", 4 },
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFilesFromHash(tc.hash)
			if tc.expected != len(actual) {
				t.Errorf("Expected %d, got %d", tc.expected, len(actual))
			}
		})
	}
}

func TestGetLsTreeOut(t *testing.T) {
	cdProjectRoot(t)
	tf := getFilesFromHash("abc9ef84782c23419bd9c61f93a1352a26f99ced")
	tcs := []struct {
		name string
		files []storage.TreeFile
		nameOnly bool
		expected string
	}{
		{ "name only", tf, true, LSTreeOutNameOnly },
		{ "full", tf, false, LSTreeOutFull },
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := getLsTreeOut(tc.files, tc.nameOnly)
			if tc.expected != actual {
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