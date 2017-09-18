#!/bin/bash
. $(dirname $0)/lib/logging

ROOT_DIR=$(realpath $(dirname $0)/..)
ARCH=x64
MAIN_GOPATH=${HOME}/go
UPLOAD_TO_PI=0

if [[ $# > 0 ]]; then
  ARCH=$1
  if [[ $# > 1 && $2 == "-c" ]]; then
    CLEAN=1
  fi

  if [[ "$@" =~ " -f" && $ARCH == "arm" ]]; then
    UPLOAD_TO_PI=1
  fi
fi

logdbg "GOPATH: ${MAIN_GOPATH}:${ROOT_DIR}"

if [[ $CLEAN == 1 ]]; then
  $(dirname $0)/clean_build.sh
fi

if [[ $ARCH == "arm" ]]; then
  loginfo "Building for linux arm (Raspberry Pi)"
  mkdir -p ${ROOT_DIR}/build/linux_arm
  GOPATH=${MAIN_GOPATH}:${ROOT_DIR}  GOOS=linux GOARCH=arm go build -o ${ROOT_DIR}/build/linux_arm/solarcmd ${ROOT_DIR}/src/cmd/solarcmd/main.go

  if [ ${UPLOAD_TO_PI} -eq 1 ]; then
    $(dirname $0)/upload_build.sh
  fi
   
elif [[ $ARCH == "x64" ]]; then
  loginfo "Building for linux x64"
  mkdir -p ${ROOT_DIR}/build/linux_amd64
#  GOPATH=${GOPATH}:${ROOT_DIR} ; cd ${ROOT_DIR}/src/cmd/solarcmd ; go install
#  mkdir -p ${ROOT_DIR}/build/arm-linux
  GOPATH=${MAIN_GOPATH}:${ROOT_DIR} go build -o ${ROOT_DIR}/build/linux_amd64/solarcmd ${ROOT_DIR}/src/cmd/solarcmd/main.go
else 
  logerr "Unsupported platform $ARCH"
  exit 1
fi

loginfo "Done."