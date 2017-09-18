#!/bin/bash
. $(dirname $0)/lib/logging

ROOT_DIR=$(realpath $(pwd)/$(dirname $0)/..)
OUTPUT_DIR=${ROOT_DIR}/build
MAIN_GOPATH=${HOME}/go

logdbg "GOPATH: ${MAIN_GOPATH}:${ROOT_DIR}"

TEST_FLAGS="-v"
if [[ $# > 0 ]]; then
  if [[ $1 == "-q" ]]; then
    TEST_FLAGS=""
  fi
fi

TESTS_FOUND=0

if [[ -f ${OUTPUT_DIR}/unit_test_report.xml ]]; then
  rm ${OUTPUT_DIR}/*unit_test_report.xml
fi

for dir in `/usr/bin/find ${ROOT_DIR}/src/cmd/solarcmd ${ROOT_DIR}/src/pkg -mindepth 1 -maxdepth 1 -type d` ; do
  if [[ "$(/usr/bin/find $dir -type f -name *_test.go)" != "" ]]; then
#  if [[ "$(ls $dir/*_test.go 2>/dev/null)" != "" ]]; then
    if [[ $TESTS_FOUND == 0 ]]; then
      loginfo "Running go unit tests..."
    fi
    TESTS_FOUND=1
    package=$(basename $dir)
    loginfo "Running unit tests for package ${package}"
    cd  $dir
    GOPATH=${MAIN_GOPATH}:${ROOT_DIR} go test ${TEST_FLAGS} | tee ${OUTPUT_DIR}/unit_test_output.txt | ${HOME}/go/bin/go-junit-report > ${OUTPUT_DIR}/${package}_unit_test_report.xml
    cat ${OUTPUT_DIR}/unit_test_output.txt
    rm ${OUTPUT_DIR}/unit_test_output.txt
    cd -
  fi
done

if [[ $TESTS_FOUND == 1 ]]; then
  loginfo "Done."
else
  logwarn "No unit tests found."
fi
