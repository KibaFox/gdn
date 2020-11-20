package gdn_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

const bufSize = 4096

func TestGdn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gdn Suite")
}

func tmpDir() string {
	tmp, err := ioutil.TempDir(".", "tmp")
	Expect(err).ToNot(HaveOccurred())

	return tmp
}

// MatchDir compares one directory with another.
func MatchDir(expected interface{}) types.GomegaMatcher {
	return &matchDir{
		expected: expected,
	}
}

type matchDir struct {
	expected interface{}
}

var ErrNotString = errors.New("type is not string")

func (matcher *matchDir) Match(actual interface{}) (success bool, err error) {
	expectedDir, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("expected %#v: %w",
			matcher.expected, ErrNotString)
	}

	Expect(expectedDir).Should(BeADirectory())

	actualDir, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("actual %#v: %w", actual, ErrNotString)
	}

	e := mapDir(expectedDir)
	a := mapDir(actualDir)

	Expect(e).ShouldNot(BeEmpty(), "expected dir should not be empty")
	Expect(a).ShouldNot(BeEmpty(), "actual dir should not be empty")
	Expect(a).Should(Equal(e))

	return true, nil
}

type pathInfo struct {
	isDir bool
	hash  string
}

// mapDir will traverse a directory and return a map.  The key is the path
// relative to the given directory. The value is the fileinfo which holds
// whether the file is a directory; if it's a file, it will also hold the file's
// sha256 hash.
func mapDir(dir string) map[string]*pathInfo {
	Expect(dir).To(BeADirectory(), "%s is not a directory", dir)

	m := make(map[string]*pathInfo)

	Expect(filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walking path %q: %w", path, err)
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			rel, err := filepath.Rel(dir, path)
			Expect(err).ToNot(HaveOccurred(),
				"could not get relative path of %s with base of %s", path, dir)

			if info.IsDir() {
				m[rel] = &pathInfo{isDir: true}
			} else {
				hash, hErr := filehash(path)
				if hErr != nil {
					return fmt.Errorf("error getting hash: %w", hErr)
				}

				m[rel] = &pathInfo{isDir: false, hash: hash}
			}

			return nil
		})).To(Succeed())

	return m
}

// filehash gets the sha256 hash of a file.  This is returned as a hex string.
func filehash(path string) (string, error) {
	file, fErr := os.Open(path)
	if fErr != nil {
		return "", fmt.Errorf("problem opening %s for reading: %w", path, fErr)
	}

	defer file.Close()

	buf, sha256 := make([]byte, bufSize), sha256.New()

	for {
		n, rErr := file.Read(buf)
		if rErr == io.EOF {
			break
		} else if rErr != nil {
			return "", fmt.Errorf("error reading %s: %w", path, rErr)
		}

		if n > 0 {
			_, shaErr := sha256.Write(buf[:n])
			if shaErr != nil {
				return "", fmt.Errorf(
					"problem writing sha256 for %s: %w", path, shaErr)
			}
		}
	}

	return hex.EncodeToString(sha256.Sum(nil)), nil
}

func (matcher *matchDir) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected directory\n\t%#v\nto be the same as\n\t%#v",
		actual, matcher.expected,
	)
}

func (matcher *matchDir) NegatedFailureMessage(actual interface{}) (message string) { // nolint: lll
	return fmt.Sprintf(
		"Expected directory\n\t%#v\nto not be the same as\n\t%#v",
		actual, matcher.expected,
	)
}
