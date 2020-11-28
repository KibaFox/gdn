package gdn_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const bufSize = 4096

// tmpDir creates a new temporary directory with the tmp prefix.
// Calls t.Fatalf() if an error occurs.
func tmpDir(t *testing.T) string {
	tmp, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Fatalf("could not create tmp dir: %v", err)
	}

	return tmp
}

// pretty converts an interface into a pretty representation by converting it to
// JSON.  This is to help debug structs by printing the values.
// Calls t.Fatalf() if an error occurs.
func pretty(t *testing.T, i interface{}) string {
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		t.Fatalf("could not pretty print %+v by converting to JSON: %v", i, err)
	}

	return string(s)
}

// pathIsRegularFile expects the given path to be a regular file.
// Calls t.Errorf() if the path is not a regular file.
func pathIsRegularFile(t *testing.T, path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("path does not exist: %s", path)
		} else {
			t.Errorf("error getting info for path: %s: %v", path, err)
		}

		return false
	}

	if !info.Mode().IsRegular() {
		t.Errorf("path is not a regular file: %s", path)
		return false
	}

	return true
}

// isDir expects the given path to be a directory.
// Calls t.Errorf() if the path is not a directory.
func pathIsDir(t *testing.T, path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("path does not exist: %s", path)
		} else {
			t.Errorf("error getting info for path: %s: %v", path, err)
		}

		return false
	}

	if !info.IsDir() {
		t.Errorf("path is not a directory: %s", path)
		return false
	}

	return true
}

// matchFile compares the contents of two files and expects them to be the same.
// Calls t.Errorf() if an they do not match.
func matchFile(t *testing.T, actualFile, expectedFile string) {
	a := filehash(t, actualFile)
	e := filehash(t, expectedFile)

	if a != e {
		t.Errorf("%s contents does not match contents of %s",
			actualFile, expectedFile)
	}
}

// matchDir compares two directories and expects them to match.
// Calls t.Errorf() if they do not match.
func matchDir(t *testing.T, actualDir, expectedDir string) {
	a := mapDir(t, actualDir)
	e := mapDir(t, expectedDir)

	for efile, einfo := range e {
		ainfo, exists := a[efile]
		if !exists {
			t.Errorf("%s not present in actual", efile)
			continue
		}

		if einfo.isDir != ainfo.isDir {
			if einfo.isDir {
				t.Errorf("%s is a directory, but actual is not", efile)
			} else {
				t.Errorf("%s is a file, but actual is not", efile)
			}

			continue
		}

		if einfo.hash != ainfo.hash {
			t.Errorf("%s contents does not match expected", efile)
		}
	}

	for afile := range a {
		if _, exists := e[afile]; !exists {
			t.Errorf("%s is present in actual, but is not in expected", afile)
		}
	}
}

type pathInfo struct {
	isDir bool
	hash  string
}

// mapDir will traverse a directory and return a map.  The key is the path
// relative to the given directory. The value is the fileinfo which holds
// whether the file is a directory; if it's a file, it will also hold the file's
// sha256 hash.
// Calls t.Fatalf() if an error occurs.
func mapDir(t *testing.T, dir string) map[string]*pathInfo {
	if !pathIsDir(t, dir) {
		t.Fatal()
	}

	m := make(map[string]*pathInfo)

	if err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking path %q: %w", path, err)
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return fmt.Errorf(
					"could not get relative path of %s with base of %s",
					path, dir)
			}

			if info.IsDir() {
				m[rel] = &pathInfo{isDir: true}
			} else {
				hash := filehash(t, path)
				m[rel] = &pathInfo{isDir: false, hash: hash}
			}

			return nil
		},
	); err != nil {
		t.Fatalf("%v", err)
	}

	return m
}

// filehash gets the sha256 hash of a file.  This is returned as a hex string.
// Calls t.Fatalf() if an error occurs.
func filehash(t *testing.T, path string) string {
	file, fErr := os.Open(path)
	if fErr != nil {
		t.Fatalf("problem opening %s for reading: %v", path, fErr)
		return ""
	}

	defer file.Close()

	buf, sha256 := make([]byte, bufSize), sha256.New()

	for {
		n, rErr := file.Read(buf)
		if errors.Is(rErr, io.EOF) {
			break
		} else if rErr != nil {
			t.Errorf("error reading %s: %v", path, rErr)
			return ""
		}

		if n > 0 {
			_, shaErr := sha256.Write(buf[:n])
			if shaErr != nil {
				t.Errorf("problem writing sha256 for %s: %v", path, shaErr)
				return ""
			}
		}
	}

	return hex.EncodeToString(sha256.Sum(nil))
}
