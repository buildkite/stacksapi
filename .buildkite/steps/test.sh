#!/usr/bin/env bash

set -Eeuo pipefail

echo '+++ Running tests'
go tool gotestsum -- -count=1 "$@" ./...
