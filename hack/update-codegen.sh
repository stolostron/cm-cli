#!/usr/bin/env bash
# Copyright Contributors to the Open Cluster Management project

set -o errexit
set -o nounset
set -o pipefail
# set -x

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/..

TMP=$(mktemp -d)
rm -rf $TMP/github.com/open-cluster-management/${PROJECT_NAME}
go install k8s.io/code-generator/cmd/{client-gen,lister-gen,informer-gen,deepcopy-gen,register-gen}

# Go installs the above commands to get installed in $GOBIN if defined, and $GOPATH/bin otherwise:
GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"

if [[ "${VERIFY_CODEGEN:-}" == "true" ]]; then
  echo "Running in verification mode"
  VERIFY_FLAG="--verify-only"
fi
COMMON_FLAGS="${VERIFY_FLAG:-} --go-header-file ${SCRIPT_ROOT}/hack/boilerplate.go.txt"

APIS_PKG=github.com/open-cluster-management/${PROJECT_NAME}/api
FQ_APIS=github.com/open-cluster-management/${PROJECT_NAME}/api/cm-cli/v1alpha1

echo "Generating register at ${FQ_APIS}"
"${gobin}/register-gen" --output-package "${FQ_APIS}" --input-dirs ${FQ_APIS} ${COMMON_FLAGS} --output-base $TMP

cp -r $TMP/github.com/open-cluster-management/${PROJECT_NAME}/* .
