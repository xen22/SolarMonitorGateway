#!/bin/bash
. $(dirname $0)/lib/logging

ROOT_DIR=$(realpath $(dirname $0)/..)

PI_HOST=pi-solar.local
if [[ $# > 0 ]]; then
  PI_HOST=$1
fi

loginfo "Uploading solarcmd (linux arm build) to ${PI_HOST}"
scp ${ROOT_DIR}/build/linux_arm/solarcmd sol@${PI_HOST}:bin/

loginfo "Uploading solarcmd configuration file to ${PI_HOST}"
scp ${ROOT_DIR}/src/cmd/solarcmd/*.config.json sol@${PI_HOST}:bin/

