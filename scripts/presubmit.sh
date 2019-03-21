#!/bin/bash

set -e

cd "${0%/*}/.."

gofmt -s -d $(find . -type f -name '*.go') |& perl -pe 'END{exit($. > 0 ? 1 : 0)}'

go test -v -race ./...

echo "Running extra checks..."
go vet ./...
staticcheck ./...
golint $(go list ./... | grep -v 'oauther') |& grep -v 'exported .* should have' | perl -pe 'END{exit($. > 0 ? 1 : 0)}'
