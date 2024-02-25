#!/bin/bash

if ! command -v golangci-lint >/dev/null 2>&1
then
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
fi

golangci-lint run
