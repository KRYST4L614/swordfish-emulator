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

## Guide for contributors

> To send pull request in repository you need to pass all CI steps that can be found in .gitlab-ci.yaml.
> To run steps locally you can find needed commands in Makefile.

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
make unit-test
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

You can see code coverage with:
```bash
make cover
```

### Run shared gitlab runner in Docker locally

You can start gitlab runner with:
```bash
make runner
```
This makes your CI available to run pipeline and pass all checks. Your PR will be ignored until CI pass.

### Generate Go structs with OpenAPI spec

For easy update to newer Redfish and Swordfish spec, we use code generation for domain models
from [OpenAPI](https://www.openapis.org/) yaml specification that can be fond in official Redfish and Swordfish resources.
See [Swordfish Schema and Registries Bundle](https://www.snia.org/sites/default/files/technical-work/swordfish/draft/v1.2.6/zip/Swordfish_v1.2.6_Schema.zip) and [Redfish Schema Bundle](https://www.dmtf.org/dsp/DSP8010).

After downloading scheme version that is used now by Swordfish API Emulator you can add new resource to generate following this steps:
1. See [existing](docs/openapi) resources and how they are included in [api.yaml](docs/openapi/api.yaml)
2. Add new resource yaml file from downloaded archive
3. Reference your resource in [api.yaml](docs/openapi/api.yaml) as in example:
```yaml
ResourceName:
  $ref: ./your-file.yaml#/components/schemas/ResourceName_ResourceName
```
Example demonstrates how to remove `Resource_Resource` naming that is practiced in Redfish and Swordfish OpenAPI specs.

4. Install needed tools:
    - [Redocly CLI](https://redocly.com/docs/cli/installation/)
    - [OpenAPI Client and Server Code Generator](https://github.com/deepmap/oapi-codegen)
5. Run `make oapi-gen`