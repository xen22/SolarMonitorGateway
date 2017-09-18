#!/bin/bash
. $(dirname $0)/lib/logging

ROOT_DIR=$(realpath $(pwd)/$(dirname $0)/..)

loginfo "Cleaning previous build"
rm -rf ${ROOT_DIR}/build/*
