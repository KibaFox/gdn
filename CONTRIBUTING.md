# Contributing

## Working with Project Code

### Dependencies

This project has the following dependencies:

* [Go](https://golang.org) the programming language and standard libraries.
* [Mage](https://magefile.org/) a build tool similar to make, but uses Go.
* [git](https://git-scm.com/) for version control.
* [golangci-lint](https://golangci-lint.run/) for linting.

You need to install these before working with the project.

It is also highly recommended to install an
[EditorConfig](https://editorconfig.org/) plugin for the editor you use.  This
will read the `.editorconfig` file for this project and use the configured
settings for your editor when you work on this project.

### Mage Targets and Help

From the root directory of the project, run `mage`:

```sh
mage
```

This will show you the list of targets made available to you via the
`magefile.go`.  If there are more details for a target you can run:

```sh
mage -h <target>
```


### Linting

The `lint` target runs static code analysis with `golangci-lint` configured via
the `.golangci.yml` for this project.

```sh
mage -v lint
```

### Testing

The `test` target runs all tests for the project.

```sh
mage -v test
```
