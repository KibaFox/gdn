//+build mage

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-bindata/go-bindata/v3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"gitlab.com/kibafox/gdn"
)

// SakuraURL is the URL to download a release of Sakura (a CSS theme).
const SakuraURL = "https://github.com/oxalorg/sakura/archive/1.3.0.tar.gz"

// Clean removes the dist/ directory and any temporary directories from testing.
func Clean() error {
	if err := os.RemoveAll("dist"); err != nil {
		return err
	}

	// Remove temporary directories.
	tmpPaths, err := filepath.Glob("tmp*")
	if err != nil {
		return err
	}

	for _, tmp := range tmpPaths {
		if err := os.RemoveAll(tmp); err != nil {
			return err
		}
	}

	return nil
}

// Build gdn into the dist/ directory.
func Build() error {
	name := "gdn"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	if err := os.MkdirAll("dist", 0777); err != nil {
		return err
	}

	return sh.Run("go", "build", "-o", fmt.Sprintf("dist/%s", name),
		"./cmd/gdn")
}

// Install will install via `go install ./cmd/gdn`.
func Install() error {
	return sh.Run("go", "install", "./cmd/gdn")
}

// Lint will perform style checks and static analysis on the Go code.
func Lint() error {
	if err := sh.Run("golangci-lint", "run"); err != nil {
		return err
	}

	return sh.Run("golangci-lint", "run", "magefile.go")
}

// Test will run all tests.
func Test() error {
	return sh.Run("ginkgo", "-v", "test", "./...")
}

// Gen is for targets that provide some form of code generation.
type Gen mg.Namespace

func (Gen) BinData() error {
	cfg := bindata.NewConfig()
	cfg.Package = "gdn"
	cfg.Input = []bindata.InputConfig{
		{
			Path:      "./assets/sakura",
			Recursive: false,
		},
	}

	if err := bindata.Translate(cfg); err != nil {
		return fmt.Errorf("error generating bindata: %w", err)
	}

	return nil
}

// Test generates testdata/expected from testdata/src.
func (Gen) Test() error {
	root := gdn.NewTree(
		filepath.Join("testdata", "src"),
		filepath.Join("testdata", "expected"),
	)

	if err := root.Scan(); err != nil {
		return err
	}

	return root.Grow()
}

type Fetch mg.Namespace

// Sakura fetches and extracts the sakura CSS theme into the ./assets/sakura
// directory.
func (Fetch) Sakura(ctx context.Context) error {
	tmp, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}

	defer os.RemoveAll(tmp)

	dl := filepath.Join(tmp, "sakura.tar.gz")

	if _, err := os.Stat(dl); err != nil {
		if err := downloadFile(ctx, dl, SakuraURL); err != nil {
			return fmt.Errorf("error downloading sakura: %w", err)
		}
	} else {
		log.Printf("%s already exists; skipping download", dl)
	}

	destDir := filepath.Join("assets", "sakura")

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("error making sakura dir: %w", err)
	}

	if err := sh.Run("tar",
		"-C", destDir,
		"--strip-components", "1",
		"-xzvf", dl,
		"*/LICENSE.txt",
		"*/css/normalize.css",
		"*/css/sakura.css",
		"*/css/sakura-dark.css",
	); err != nil {
		return fmt.Errorf("error extracting sakura: %w", err)
	}

	var (
		srcpath = filepath.Join(destDir, "source.txt")
		srctxt  = []byte(SakuraURL + "\n")
	)

	if err := ioutil.WriteFile(srcpath, srctxt, 0600); err != nil {
		return fmt.Errorf("error writing download URL to %s: %w", srcpath, err)
	}

	return nil
}

// downloadFile will download a url to a local file.
func downloadFile(ctx context.Context, dst, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting %s: %w", url, err)
	}

	defer resp.Body.Close()

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", dst, err)
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying response to %s: %w", dst, err)
	}

	return nil
}
