# Muffet

[![GitHub Action](https://img.shields.io/github/actions/workflow/status/raviqqe/muffet/test.yaml?branch=main&style=flat-square)](https://github.com/raviqqe/muffet/actions)
[![Codecov](https://img.shields.io/codecov/c/github/raviqqe/muffet.svg?style=flat-square)](https://codecov.io/gh/raviqqe/muffet)
[![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/muffet?style=flat-square)](https://goreportcard.com/report/github.com/raviqqe/muffet)
[![Docker](https://img.shields.io/docker/pulls/raviqqe/muffet?style=flat-square)](https://hub.docker.com/r/raviqqe/muffet)
[![License](https://img.shields.io/github/license/raviqqe/muffet.svg?style=flat-square)](LICENSE)

![demo](img/demo.gif)

Muffet is a website link checker which scrapes and inspects all pages in a
website recursively.

## Features

- Massive speed
- High compatibility with web browsers
- Different tag support (`a`, `img`, `link`, `script`, etc)
- Multiple output formats (text, JSON, and JUnit XML)

## Installation

```sh
go install github.com/raviqqe/muffet/v2@latest
```

### Homebrew

```sh
brew install muffet
```

## Usage

```sh
muffet https://shady.bakery.hotland
```

For more information, see `muffet --help`.

### Docker

```sh
docker run raviqqe/muffet https://shady.bakery.hotland
```

### GitHub Action

- [My Broken Link Checker](https://github.com/ruzickap/action-my-broken-link-checker)

Currently, we do not provide any official one. Feel free to create an issue if you want!

## License

[MIT](LICENSE)
