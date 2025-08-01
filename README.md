<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/fz"><img src="https://kaos.sh/r/fz.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/y/fz"><img src="https://kaos.sh/y/cb41fe8e630a408c86a8c227393f5359.svg" alt="Codacy Badge" /></a>
  <a href="https://kaos.sh/w/fz/ci"><img src="https://kaos.sh/w/fz/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/fz/codeql"><img src="https://kaos.sh/w/fz/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#usage-demo">Demo</a> • <a href="#installation">Installation</a> • <a href="#upgrading">Upgrading</a> • <a href="#command-line-completion">Completions</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#license">License</a></p>

<br/>

`fz` is a simple tool for formatting [`go-fuzz`](https://github.com/dvyukov/go-fuzz) output.

### Usage demo

[![demo](https://github.com/user-attachments/assets/0f350911-3c63-4e02-8eec-dc038c034350)](#usage-demo)

### Installation

#### From source

To build the `fz` from scratch, make sure you have a working [Go 1.23+](https://github.com/essentialkaos/.github/blob/master/GO-VERSION-SUPPORT.md) workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/fz@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/fz/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) fz
```

### Upgrading

Since version `1.2.0` you can update `fz` to the latest release using [self-update feature](https://github.com/essentialkaos/.github/blob/master/APPS-UPDATE.md):

```bash
fz --update
```

This command will runs a self-update in interactive mode. If you want to run a quiet update (_no output_), use the following command:

```bash
fz --update=quiet
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo fz --completion=bash 1> /etc/bash_completion.d/fz
```


ZSH:
```bash
sudo fz --completion=zsh 1> /usr/share/zsh/site-functions/fz
```


Fish:
```bash
sudo fz --completion=fish 1> /usr/share/fish/vendor_completions.d/fz.fish
```

### Man documentation

You can generate man page for `fz` using next command:

```bash
fz --generate-man | sudo gzip > /usr/share/man/man1/fz.1.gz
```

### Usage

<img src=".github/images/usage.svg" />

### CI Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/fz/ci.svg?branch=master)](https://kaos.sh/w/fz/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/fz/ci.svg?branch=develop)](https://kaos.sh/w/fz/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
