#!/bin/bash

set -ex

go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

if [ "$CI" == "on" ]; then
  include_cov=coverage.txt bash <(curl -s https://codecov.io/bash) -t "$CODECOV_TOKEN"
fi
