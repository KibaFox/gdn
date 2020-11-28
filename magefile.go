// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Clean removes the ./dist/ directory and any temp directories.
func Clean() error {
	if err := os.RemoveAll("dist"); err != nil {
		return fmt.Errorf("could not remove ./dist/ dir: %w", err)
	}

	// Remove temporary directories.
	tmpPaths, err := filepath.Glob("tmp*")
	if err != nil {
		return fmt.Errorf("error globbing `tmp*`: %w", err)
	}

	for _, tmp := range tmpPaths {
		if err := os.RemoveAll(tmp); err != nil {
			return fmt.Errorf("could not remove tmp dir: %s: %w", tmp, err)
		}
	}

	return nil
}

// Build compiles `gdn` into the ./dist/ directory.
func Build() error {
	if err := os.MkdirAll("dist", 0o755); err != nil {
		return fmt.Errorf("could not make ./dist/ dir: %w", err)
	}

	if err := run("go", "build", "-o", "dist", "./cmd/gdn"); err != nil {
		return fmt.Errorf("problem building `gdn`: %w", err)
	}

	return nil
}

// Install will install `gdn` via go install.
func Install() error {
	if err := run("go", "install", "./cmd/gdn"); err != nil {
		return fmt.Errorf("problem installing `gdn`: %w", err)
	}

	return nil
}

// Lint runs golangci-lint on the project.
func Lint() error {
	if err := run("golangci-lint", "run"); err != nil {
		return fmt.Errorf("problem linting project: %w", err)
	}

	if err := run("golangci-lint", "run", "magefile.go"); err != nil {
		return fmt.Errorf("problem linting magefile.go: %w", err)
	}

	return nil
}

// Test runs all tests for the project.
func Test() error {
	if err := run("go", "test", "./..."); err != nil {
		return fmt.Errorf("problem testing project: %w", err)
	}

	return nil
}

func run(cmd string, args ...string) (err error) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	fmt.Println("exec:", cmd, strings.Join(args, " "))

	return c.Run()
}
