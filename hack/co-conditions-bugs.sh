#!/usr/bin/env bash

set -euo pipefail

git_commit="$( git describe --tags --always --dirty )"
build_date="$( date -u '+%Y%m%d' )"
version="v${build_date}-${git_commit}"

OUTPUT_DIR="./output"
OUTPUT_FILE="${OUTPUT_DIR}/co-conditions-bugs.md"

mkdir -p "${OUTPUT_DIR}"
go run ./cmd/co-conditions-bugs --output-format md --output-file "${OUTPUT_FILE}"

ORIGIN_DIR="../../openshift/origin"

git_commit_origin="$( git -C ${ORIGIN_DIR} describe --tags --always --dirty )"

cat << EOF >> "${OUTPUT_FILE}"

## Build info

* Build: ${version}

* \`git_commit_origin\`: ${git_commit_origin}
EOF