<p align="center"><a href="#readme"><img src="https://gh.kaos.st/fz.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/fz"><img src="https://travis-ci.com/essentialkaos/fz.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/fz"><img src="https://goreportcard.com/badge/github.com/essentialkaos/fz"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-fz-master"><img alt="codebeat badge" src="https://codebeat.co/badges/64a79279-c198-422c-862c-d4e735358ac1" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

`fz` is a simple tool for formatting [`go-fuzz`](https://github.com/dvyukov/go-fuzz) output.

### Usage demo

[![demo](https://gh.kaos.st/fz-001.gif)](#usage-demo)

### Installation

#### From source

Before the initial install, allow git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

To build the `fz` from scratch, make sure you have a working Go 1.10+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/fz
```

If you want to update `fz` to latest stable release, do:

```
go get -u github.com/essentialkaos/fz
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/fz/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) fz
```

### Usage

```
Usage: go-fuzz … |& fz

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/fz.svg?branch=master)](https://travis-ci.com/essentialkaos/fz) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/fz.svg?branch=develop)](https://travis-ci.com/essentialkaos/fz) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
