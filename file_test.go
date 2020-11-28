package gdn_test

import (
	"os"
	"path/filepath"
	"testing"

	"git.sr.ht/~kiba/gdn"
)

func TestTypeByExtension(t *testing.T) {
	tbls := []struct {
		ext      string
		expected gdn.FileType
	}{
		{".md", gdn.Markdown},
		{".mkd", gdn.Markdown},
		{".markdown", gdn.Markdown},
		{".jpeg", gdn.Unknown},
		{".txt", gdn.Unknown},
		{".unknown", gdn.Unknown},
	}

	for _, tbl := range tbls {
		result := gdn.TypeByExtension(tbl.ext)
		if result != tbl.expected {
			t.Errorf("TypeByExtension(%s) gave: %s, expecting: %s",
				tbl.ext, result, tbl.expected)
		}
	}
}

func TestFileTypeString(t *testing.T) {
	tbls := []struct {
		ftype    gdn.FileType
		expected string
	}{
		{gdn.Markdown, "Markdown"},
		{gdn.Unknown, "Unknown"},
		{gdn.FileType(256), "Unknown"},
	}

	for _, tbl := range tbls {
		result := tbl.ftype.String()
		if result != tbl.expected {
			t.Errorf("FileType(%d).String() gave: %s, expecting: %s",
				tbl.ftype, result, tbl.expected)
		}
	}
}

func TestCopyFile(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)

	src := filepath.Join("testdata", "src", "example", "mytext.txt")
	dst := filepath.Join(tmp, "mytext.txt")

	if err := gdn.CopyFile(src, dst); err != nil {
		t.Fatalf("error copying source and destination: %v", err)
	}

	pathIsRegularFile(t, dst)
	matchFile(t, src, dst)
}

func TestChExt(t *testing.T) {
	tbls := []struct {
		path     string
		ext      string
		expected string
	}{
		{"some/file", ".html", "some/file.html"},
		{"/asdf/qwer/markdown.md", ".html", "/asdf/qwer/markdown.html"},
		{"some/image.png", ".jpeg", "some/image.jpeg"},
	}

	for _, tbl := range tbls {
		result := gdn.ChExt(tbl.path, tbl.ext)
		if result != tbl.expected {
			t.Errorf("ChExt(%s, %s) gave: %s, expecting: %s",
				tbl.path, tbl.ext, result, tbl.expected)
		}
	}
}
