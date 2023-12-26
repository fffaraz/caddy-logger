#!/bin/bash
set -euxo pipefail

go get -u ./...
go mod tidy
~/go/bin/gofumpt -l -w -extra .
