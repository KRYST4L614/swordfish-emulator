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
    - [Fast start](#fast-start)
    - [Configuration](#configuration)
  - [Guide for contributors](#guide-for-contributors)
    - [Download project](#download-project)
    - [Fetch Go dependencies](#fetch-go-dependencies)
    - [Run unit test and update coverage report](#run-unit-test-and-update-coverage-report)
    - [Generate mocks for testing](#generate-mocks-for-testing)
    - [Run linter](#run-linter)
    - [See coverage](#see-coverage)
    - [Run shared gitlab runner in Docker locally](#run-shared-gitlab-runner-in-docker-locally)
    - [Generate Go structs with OpenAPI spec](#generate-go-structs-with-openapi-spec)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Swordfish API Emulator

## About the project

This project provides server application that emulates Swordfish API.

> As Swordfish specification extends Redfish our emulator supports both of them

### Docs

- [Technical requirements, development requirements and design](docs/design/[Swordfish%20API%20Emulator]%20Документация.pdf)
- [Swordfish API documentation](https://www.snia.org/sites/default/files/technical-work/swordfish/release/v1.2.5a/html/Specification/Swordfish_v1.2.5a_Specification.html)
- [Redfish API documentation](https://www.dmtf.org/sites/default/files/standards/documents/DSP0266_1.20.0.html)

### Status

The project is currently under development.

### See also

- [Ansible-collection for Swordfish API](https://gitlab.com/IgorNikiforov/swordfish_ansible_plugin)

## User guide

### Fast start

To run emulator just download repository with command
```bash
git clone https://gitlab.com/IgorNikiforov/swordfish-emulator-go.git
```
and start server with command
```bash
make start
```
or explicitly
```bash
go run cmd/emulator/main.go
```

Then server is started and you can access all Redfish/Swordfish endpoints from
service root endpoint `http://localhost:<port from configuration>/redfish/v1`

> Default port is `8080`

### Configuration

Our emulator provides broad configuration opportunities, let's cover main options.
Main configuration file is placed [here](configs/emulator/config.yaml). There you can find main settings:
- `db`: configures connection to database, if you don't need your own PostgreSQL instance - don't change this. We manage this setting
by ourselves for appropriate connection with embedded PostgreSQL.
- `embedded-psql`: configures embedded instance of PostgreSQL, so you don't need to think about invoking your own
- `dataset`:
  - `path`: with this setting you can choose dataset, that contains basic Swordfish scheme. This scheme will be loaded,
when emulator start. You can provide your own scheme and pass relative path here, but you dataset should be described as it is done in default one, otherwise it can be loaded with errors
  - `overwrite`: this option says if emulator need to overwrite your last session progress, persisted in database, or save it
- `server`: here you can find default server settings

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