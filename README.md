[![Build Status](https://gitlab.com/IgorNikiforov/swordfish-emulator-go/badges/main/pipeline.svg?key_text=CI)](https://gitlab.com/IgorNikiforov/swordfish-emulator-go/commits/main)
[![Coverage](https://gitlab.com/IgorNikiforov/swordfish-emulator-go/badges/main/coverage.svg?key_text=Coverage)](https://gitlab.com/IgorNikiforov/swordfish-emulator-go/-/commits/main)
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
<!--   *generated with [DocToc](https://github.com/thlorenz/doctoc)* -->

- [Swordfish API Emulator](#swordfish-api-emulator)
  - [About the project](#about-the-project)
    - [Docs](#docs)
    - [Status](#status)
    - [See also](#see-also)
  - [Getting started](#getting-started)
    - [Download project](#download-project)
    - [Docker deployment](#docker-deployment)
    - [Run unit test and update coverage badge](#run-unit-test-and-update-coverage-badge)
    - [Build project](#build-project)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Swordfish API Emulator

## About the project

This project provides server application that emulates Swordfish API.

### Docs

- [Technical requirements, development requirements and design](docs/design/[Swordfish%20API%20Emulator]%20Документация.pdf)
- [API documentation](api/)

### Status

The project is currently under development.

### See also

- [Ansible-collection for Swordfish API](https://gitlab.com/IgorNikiforov/swordfish_ansible_plugin)

## Getting started

### Download project

To download use:
```bash
git clone https://gitlab.com/IgorNikiforov/swordfish-emulator-go.git
```

### Docker deployment

You can use docker deployment, to prepare docker images use:
> You should have installed **Docker** and **docker-compose** on your machine

```bash
make build-images
```

Then you can simply run
```bash
make start-dc
```
to run docker container using docker-compose, and

```bash
make stop-dc
```
to stop docker containers

### Run unit test and update coverage badge

This will **run unit tests** and **update link** for coverage badge in README
```bash
make unit-test
```
You can regenerate mocks with
```bash
make gogen
```

### Build project

You can build executable files with
```bash
make build
```
Executables will be in **/bin** folder
