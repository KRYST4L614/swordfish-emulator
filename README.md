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
  - [User guide](#user-guide)
  - [Getting started with development](#getting-started-with-development)
    - [Download project](#download-project)
    - [Fetch Go dependencies](#fetch-go-dependencies)
    - [Run unit test and update coverage report](#run-unit-test-and-update-coverage-report)
    - [Generate mocks for testing](#generate-mocks-for-testing)
    - [Run linter](#run-linter)
    - [See coverage](#see-coverage)

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

## User guide

The project is currently under development, so you can't use it now :(

## Getting started with development

### Download project

To download use:
```bash
git clone https://gitlab.com/IgorNikiforov/swordfish-emulator-go.git
```

### Fetch Go dependencies

To **get dependencies**
```bash
make dep
```

### Run unit test and update coverage report

This will **run unit tests** and generate coverage report in **coverage.out** file
```bash
make test-unit
```

### Generate mocks for testing

You can regenerate mocks with
```bash
make gogen
```

### Run linter

You can run linter check with:
```bash
make gogen
```

### See coverage

You can see code coverage with
```bash
make cover
```
