<p align="center"><a href="#readme"><img src="https://gh.kaos.st/fz.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/fz/ci"><img src="https://kaos.sh/w/fz/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/fz/codeql"><img src="https://kaos.sh/w/fz/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/r/fz"><img src="https://kaos.sh/r/fz.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/fz"><img src="https://kaos.sh/b/64a79279-c198-422c-862c-d4e735358ac1.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`fz` is a simple tool for formatting [`go-fuzz`](https://github.com/dvyukov/go-fuzz) output.

### Usage demo

[![demo](https://gh.kaos.st/fz-003.gif)](#usage-demo)

### Installation

#### From source

To build the `fz` from scratch, make sure you have a working Go 1.16+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go install github.com/essentialkaos/fz
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
| `master` | [![CI](https://kaos.sh/w/fz/ci.svg?branch=master)](https://kaos.sh/w/fz/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/fz/ci.svg?branch=develop)](https://kaos.sh/w/fz/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
