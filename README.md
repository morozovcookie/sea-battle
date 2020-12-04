sea-battle
==========

Test task for OZON.ru which emulate "Sea Battle" game behaviour. [Full description](DESCRIPTION.md).


# Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Build](#build)
    - [Console](#build-with-console)
    - [Docker](#build-with-docker)
    - [Werf](#build-with-werf)
- [Usage](#usage)
    - [Console](#run-with-console)
    - [Docker](#run-with-docker)
    - [Werf](#run-with-werf)


# Requirements

- GoLang >= 1.15.3
- Docker >= 19.03.13
- Werf >= v1.1.21+fix32


# Installation

```bash
# with go get
$ go get github.com/morozovcookie/sea-battle/cmd/sea-battle/...

# with git
$ git clone https://github.com/morozovcookie/sea-battle.git
```


# Build

## Build With Console

```bash
# build binary file
$ make go-build
```

## Build With Docker

```bash
# build docker image
$ make docker-build

# publish
$ make docker-publish DOCKER_REPOSITORY=sample.registry.com
```

## Build With Werf

```bash
# build docker image
$ make werf-build

# publish docker image
$ make werf-publish DOCKER_REPOSITORY=sample.registry.com
```


# Usage

## Run With Console

```bash
$ SERVER_ADDRESS=<address> make go-run
```

## Run With Docker

```bash
$ make docker-run
```

## Run With Werf

```bash
$ make werf-run
```
