# gdn

[![builds.sr.ht status](https://builds.sr.ht/~kiba/gdn.svg)](https://builds.sr.ht/~kiba/gdn?)

Gemini digital garden static site generator.

A digital garden is like a personal knowledge-base or wiki.  Rather than a blog
where you write articles that are tied to a moment in time, a digital garden
encourages revisiting pages and changing them as you grow.

`gdn` expects your digital garden to be written in Gemini text files (with file
extension `.gmi` or `.gemini`) organized in folders.  `gdn` renders the digital
Garden into a static HTML site that can be served over HTTP/HTTPS by most web
servers.

Gemini is both a network protocol and a text format.  `gdn` renders your digital
garden to HTML so it can be viewed by a normal web browser on the web.

- <https://gemini.circumlunar.space/>

`gdn` renders Gemini text formatted files into HTML to produce a static site.

## Work in Progress

This is currently a work in progress.  Not all features work.

## Building and Installing from Source Code

### Dependencies

This project has the following dependencies:

* [Go](https://golang.org) the programming language and standard libraries.
* [Mage](https://magefile.org/) a build tool similar to make, but uses Go.
* [git](https://git-scm.com/) for version control.

You need to install these before building or installing.

### Fetching Code

Clone this repository using `git` and change directory:

```sh
git clone https://git.sr.ht/~kiba/gdn
cd gdn
```

### Building from Source

The `build` target compiles from source and places the executable in the
`./dist/` directory.

```sh
mage -v build
```

### Installing from Source

The `install` target compiles from source and places the executable into your
`$GOPATH/bin` directory via `go install`.  This must be in your `PATH` if you
wish to execute `gdn` without specifying the full path of `gdn`.

```sh
mage -v install
```

## Contributing

Contributions are welcome.  Please read the [CONTRIBUTING.md](CONTRIBUTING.md)
guide.

## License

This project is licensed by the MIT License.  See: [LICENSE](LICENSE).
