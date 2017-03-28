# s3
A general command line client for AWS S3.

## Install

```
# Set a directory where Go packages will be installed.
export GOPATH=$HOME/go

# Install `s3` and dependencies.
go get github.com/xstevens/s3
```

This will make `s3` available as `$GOPATH/bin/s3`.
Add `$GOPATH/bin` to your `$PATH` to have it available as `s3`.

### Build

To run subsequent builds, use `go build`:

```
# Ensure you're in the `s3` source directory.
cd $GOPATH/src/github.com/xstevens/s3

# Run the build.
go build
```

### Cross-compiling

With Go 1.5 or above, cross-compilation support is built in.
See [Dave Cheney's blog post](http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5)
for a tutorial and the [golang.org docs](https://golang.org/doc/install/source#environment)
for details on `GOOS` and `GOARCH` values for various target operating systems.

A typical build for Linux would be:
```
# Ensure you're in the `s3` source directory.
cd $GOPATH/src/github.com/xstevens/s3

# Run the build.
GOOS=linux GOARCH=amd64 go build -o s3_linux_amd64
```

## Usage
```
$ ./s3 -h
usage: s3 [<flags>] <command> [<args> ...]

A general command line client for S3.

Flags:
  -h, --help     Show context-sensitive help (also try --help-long and --help-man).
      --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  cat --bucket=BUCKET --prefix=PREFIX [<flags>]
    Reads all keys with specified prefix and writes them to stdout.

  upload --bucket=BUCKET --prefix=PREFIX --sourcedir=SOURCEDIR [<flags>]
    Upload file(s) to S3

  meta --bucket=BUCKET --prefix=PREFIX [<flags>]
    Reads all keys with specified prefix and writes their metadata to stdout.
```

## License
All aspects of this software are distributed under the MIT License. See LICENSE file for full license text.
