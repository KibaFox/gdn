package main

import (
	"log"
	"os"

	"gitlab.com/kibafox/gdn"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	root := gdn.NewTree(dir, "dist")

	if err := root.Scan(); err != nil {
		log.Fatal(err)
	}

	if err := root.Grow(); err != nil {
		log.Fatal(err)
	}
}
