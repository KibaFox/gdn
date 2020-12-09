package gdn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var (
	// ErrSrcNotSet occurs when the source path is not set.
	ErrSrcNotSet = errors.New("source path is not set")
	// ErrDstNotSet occurs when the destinatio path is not set.
	ErrDstNotSet = errors.New("destination path is not set")
	// ErrNotScanned occurs when Scan was not called before Grow.
	ErrNotScanned = errors.New("need to scan tree before growing it")
	// ErrEmptyTree occurs when Scan results in an empty tree.
	ErrEmptyTree = errors.New("scan resulted in an empty tree")
)

// Branch represents a directory tree used to generate the pages.
type Branch struct {
	Src      string
	Dst      string
	Path     string
	Branches []*Branch
	Leaves   []*Leaf
}

// NewTree creates the root of the tree.  The input path is the path with all
// the source files and directories (e.g. Markdown, images, etc).  The output
// path is the destination for the generated site.  This returns a single Branch
// with no Leaves which an be used to Scan the input path to populate the tree.
func NewTree(inputPath, outputPath string) Branch {
	return Branch{
		Src:  inputPath,
		Dst:  outputPath,
		Path: "/",
	}
}

// Scan will scan the input path for items to generate the site and build the
// tree.  Directories are added as Branches. Files are added as Leaves.
// Hidden files and directories are ignored.
func (b *Branch) Scan() error {
	if b.Src == "" {
		return ErrSrcNotSet
	}

	if b.Dst == "" {
		return ErrDstNotSet
	}

	files, err := ioutil.ReadDir(b.Src)
	if err != nil {
		return fmt.Errorf("could not scan directory: %s: %w", b.Src, err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			// Skip hidden directories and files.
			continue
		}

		if f.IsDir() {
			branch := &Branch{
				Src:  filepath.Join(b.Src, f.Name()),
				Dst:  filepath.Join(b.Dst, f.Name()),
				Path: filepath.Join(b.Path, f.Name()),
			}

			err := branch.Scan()
			if errors.Is(err, ErrEmptyTree) {
				// Skip empty branches
				continue
			} else if err != nil {
				return err
			}

			b.Branches = append(b.Branches, branch)
		} else {
			b.Leaves = append(b.Leaves, &Leaf{
				Src:    filepath.Join(b.Src, f.Name()),
				DstDir: b.Dst,
				Path:   filepath.Join(b.Path, f.Name()),
				Typ:    TypeByExtension(filepath.Ext(f.Name())),
			})
		}
	}

	if len(b.Leaves) == 0 && len(b.Branches) == 0 {
		return ErrEmptyTree
	}

	return nil
}

// BranchPerm sets the permission for the directories produced when growing.
const BranchPerm os.FileMode = 0750

// Grow generates the site from the branch.
func (b Branch) Grow() error {
	if b.Src == "" {
		return ErrSrcNotSet
	}

	if b.Dst == "" {
		return ErrDstNotSet
	}

	if len(b.Leaves) == 0 && len(b.Branches) == 0 {
		return ErrNotScanned
	}

	if err := os.MkdirAll(b.Dst, BranchPerm); err != nil {
		return fmt.Errorf("error making directory: %s: %w", b.Dst, err)
	}

	for _, leaf := range b.Leaves {
		if err := leaf.Grow(); err != nil {
			return err
		}
	}

	for _, branch := range b.Branches {
		if err := branch.Grow(); err != nil {
			return err
		}
	}

	return nil
}

// Leaf represnts a file.  If it is a Markdown file it will be generated into a
// page.
type Leaf struct {
	Src    string
	DstDir string
	Path   string
	Typ    FileType
}

// LeafPerm is the permission to set for the generated file the leaf produces.
const LeafPerm os.FileMode = 0640

// Dst is the destination file path for the leaf when Grow is executed.
func (l Leaf) Dst() string {
	switch l.Typ {
	case Markdown:
		return ChExt(filepath.Join(l.DstDir, filepath.Base(l.Src)), ".html")
	case Unknown:
		return filepath.Join(l.DstDir, filepath.Base(l.Src))
	default:
		return filepath.Join(l.DstDir, filepath.Base(l.Src))
	}
}

// Grow will generate a page for the leaf.
func (l Leaf) Grow() error {
	if l.Src == "" {
		return ErrSrcNotSet
	}

	if l.DstDir == "" {
		return ErrDstNotSet
	}

	switch l.Typ {
	case Markdown:
		m, readErr := ioutil.ReadFile(l.Src)
		if readErr != nil {
			return fmt.Errorf("error reading %s: %w", l.Src, readErr)
		}

		writeErr := ioutil.WriteFile(l.Dst(), blackfriday.Run(m), LeafPerm)
		if writeErr != nil {
			return fmt.Errorf("error writing %s: %w", l.Dst(), readErr)
		}

	case Unknown:
		err := CopyFile(l.Src, l.Dst())
		if err != nil {
			return fmt.Errorf("error copying %s to %s: %w", l.Src, l.Dst(), err)
		}

	default:
		err := CopyFile(l.Src, l.Dst())
		if err != nil {
			return fmt.Errorf("error copying %s to %s: %w", l.Src, l.Dst(), err)
		}
	}

	return nil
}
