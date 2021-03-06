package gdn

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileType represents a high level file type.
type FileType uint

const (
	// Unknown is an unknown file type.
	Unknown FileType = iota
	// Markdown is a markdown file type.
	Markdown
)

// String returns the string representation of FileType.  For example, if the
// FileType is Markdown it will return the string "Markdown".
func (t FileType) String() string {
	switch t {
	case Markdown:
		return "Markdown"
	case Unknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// KnownFileTypes is a map of extensions to their FileType.
var KnownFileTypes = map[string]FileType{ // nolint: gochecknoglobals
	".md":       Markdown,
	".mkd":      Markdown,
	".markdown": Markdown,
}

// TypeByExtension will look up the type by its extension.
func TypeByExtension(ext string) FileType {
	typ, ok := KnownFileTypes[ext]
	if !ok {
		return Unknown
	}

	return typ
}

// ChExt takes a path and replaces any file extension that path has with the
// the given new extension.
func ChExt(path, ext string) string {
	e := filepath.Ext(path)
	p := strings.TrimSuffix(path, e)

	return p + ext
}

// CopyFile will copy a file from the given source to the destination.
func CopyFile(src, dest string) error {
	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open src (%s) to copy: %w", src, err)
	}
	defer input.Close()

	output, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create dest (%s) to copy: %w", dest, err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		return fmt.Errorf("error copying (%s) to (%s): %w", src, dest, err)
	}

	return output.Close()
}
