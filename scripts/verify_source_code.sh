#!/bin/bash
. $(dirname $0)/lib/logging

ROOT_DIR=$(realpath $(dirname $0)/..)

loginfo "Running \"go vet\" and \"golint\" on the source code..."

for go_file in `/usr/bin/find ${ROOT_DIR}/src/cmd/solarcmd ${ROOT_DIR}/src/pkg -type f -name *.go` ; do
  if [[ "$go_file" != *"_old"* && "$go_file" != *"_vendor"* ]]; then

    vet_out=$(go vet $go_file 2>&1)
    IFS=$'\n' lines=($vet_out)
    for line in "${lines[@]}"; do
      echo "go-vet warning: $line"
    done

    lint_out=$(golint $go_file 2>&1)
    IFS=$'\n' lines=($lint_out)
    for line in "${lines[@]}"; do
      echo "go-lint warning: $line"
    done

  fi
done


