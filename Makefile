.PHONY: build clean help install lint test

help:
	@echo 'Usage: make [TARGET...]'
	@echo ''
	@echo 'Targets:'
	@echo '  clean      removes the ./dist/ directory and any temp directories'
	@echo '  build      compiles `gdn` into the ./dist/ directory'
	@echo '  install    will install `gdn` via go install'
	@echo '  lint       runs golangci-lint on the project'
	@echo '  test       runs all tests for the project'

clean:
	rm -rf dist

build:
	mkdir -p dist
	go build -o dist ./cmd/gdn

install:
	go install ./cmd/gdn

lint:
	golangci-lint run

test:
	go test ./...
