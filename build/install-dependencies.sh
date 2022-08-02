#!/bin/bash -e

# Copyright Contributors to the Open Cluster Management project

# Go tools
_OS=$(go env GOOS)
_ARCH=$(go env GOARCH)

if ! which patter > /dev/null;     then echo "Installing patter ..."; go install github.com/apg/patter@bd185be70ac8aa766084cfa431f6bfbc266795d4; fi
if ! which gocovmerge > /dev/null; then echo "Installing gocovmerge..."; go install github.com/wadey/gocovmerge@b5bfa59ec0adc420475f97f89b58045c721d761c; fi

# Build tools

# Image tools

# Check tools
