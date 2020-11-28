package gdn_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"git.sr.ht/~kiba/gdn"
)

const testsrc = "testdata/src"

func TestNewTree(t *testing.T) {
	root := gdn.NewTree("mygarden/src", "mygarden/dest")
	expected := gdn.Branch{
		Src:  "mygarden/src",
		Dst:  "mygarden/dest",
		Path: "/",
	}

	if !reflect.DeepEqual(root, expected) {
		t.Errorf("new tree %+v does not match expected %+v", root, expected)
	}
}

func TestBranchScan(t *testing.T) {
	t.Logf("+test scanning %s", testsrc)

	root := gdn.NewTree(testsrc, "tmp")

	if err := root.Scan(); err != nil {
		t.Fatalf("scan encountered an unexpected error: %v", err)
	}

	expected := gdn.Branch{
		Src:  testsrc,
		Dst:  "tmp",
		Path: "/",
		Branches: []*gdn.Branch{
			{
				Src:  testsrc + "/example",
				Dst:  "tmp/example",
				Path: "/example",
				Leaves: []*gdn.Leaf{
					{
						Src:    testsrc + "/example/mydoc.md",
						DstDir: "tmp/example",
						Path:   "/example/mydoc.md",
						Typ:    gdn.Markdown,
					},
					{
						Src:    testsrc + "/example/mytext.txt",
						DstDir: "tmp/example",
						Path:   "/example/mytext.txt",
						Typ:    gdn.Unknown,
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(root, expected) {
		t.Errorf("scanned tree %s does not match expected %s",
			pretty(t, root), pretty(t, expected))
	}

	t.Log("-test when source path is not set")

	root = gdn.NewTree("", "tmp")
	if err := root.Scan(); !errors.Is(err, gdn.ErrSrcNotSet) {
		t.Error("expected ErrSrcNotSet when source path is not set")
	}

	t.Log("-test when destination path is not set")

	root = gdn.NewTree(testsrc, "")
	if err := root.Scan(); !errors.Is(err, gdn.ErrDstNotSet) {
		t.Error("expected ErrDstNotSet when destination path is not set")
	}

	t.Log("-test when scanning an empty directory")

	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)

	root = gdn.NewTree(tmp, "tmp")
	if err := root.Scan(); !errors.Is(err, gdn.ErrEmptyTree) {
		t.Error("expected ErrEmptyTree when scanning an empty directory")
	}
}

func TestBranchGrow(t *testing.T) {
	tmp := tmpDir(t)
	defer os.RemoveAll(tmp)

	t.Log("+test that growing testdata/src matches testdata/expected")

	root := gdn.NewTree(testsrc, tmp)

	if err := root.Scan(); err != nil {
		t.Fatalf("scan encountered an unexpected error: %v", err)
	}

	if err := root.Grow(); err != nil {
		t.Fatalf("grow encountered an unexpected error: %v", err)
	}

	matchDir(t, tmp, "testdata/expected")

	t.Log("-test when the source path is not set")

	root = gdn.NewTree("", "tmp")
	if err := root.Grow(); !errors.Is(err, gdn.ErrSrcNotSet) {
		t.Error("expected ErrSrcNotSet when source path is not set")
	}

	t.Log("-test when destination path is not set")

	root = gdn.NewTree(testsrc, "")
	if err := root.Grow(); !errors.Is(err, gdn.ErrDstNotSet) {
		t.Error("expected ErrDstNotSet when destination path is not set")
	}

	t.Log("-test to warn with the branch is empty and not scanned")

	root = gdn.NewTree(testsrc, "tmp")
	if err := root.Grow(); !errors.Is(err, gdn.ErrNotScanned) {
		t.Error("expected ErrNotScanned when scan was not run")
	}
}

func TestLeafDst(t *testing.T) {
	tbls := []struct {
		leaf     gdn.Leaf
		expected string
	}{
		{
			gdn.Leaf{
				Src:    "asdf/my.md",
				DstDir: "qwer",
				Path:   "/my.md",
				Typ:    gdn.Markdown,
			},
			"qwer/my.html",
		},
		{
			gdn.Leaf{
				Src:    "asdf/my.txt",
				DstDir: "qwer",
				Path:   "/my.txt",
				Typ:    gdn.Unknown,
			},
			"qwer/my.txt",
		},
	}

	for _, tbl := range tbls {
		result := tbl.leaf.Dst()
		if result != tbl.expected {
			t.Errorf("Leaf %+v .DST() gave: %s, expecting: %s",
				tbl.leaf, result, tbl.expected)
		}
	}
}

func TestLeafGrow(t *testing.T) {
	t.Log("-test ensures error is given when source path is not set")

	leaf := gdn.Leaf{Src: "", DstDir: "tmp"}
	if err := leaf.Grow(); !errors.Is(err, gdn.ErrSrcNotSet) {
		t.Error("expected ErrSrcNotSet when source path is not set")
	}

	t.Log("-test ensures error is given when destination path is not set")

	leaf = gdn.Leaf{Src: "tmp", DstDir: ""}
	if err := leaf.Grow(); !errors.Is(err, gdn.ErrDstNotSet) {
		t.Error("expected ErrDstNotSet when destination path is not set")
	}
}
