#!/bin/bash
. $(dirname $0)/lib/logging

PLATFORM="linux_amd64"
EXPECTED_VERSION="0.0.0.1"
PI_TEST_SYSTEM="pi-solar.local"
if [[ $# -lt 1 ]]; then
  loginfo "Usage: $0 [<version>]"
  loginfo "Attempting to retrieve version from the most recent git tag on the current branch."
  EXPECTED_VERSION=$(git tag -l --points-at HEAD | tail -1)
  if [[ $EXPECTED_VERSION == "" ]]; then 
    logerr "Could not retrieve version."
    exit 1
  fi
else
  EXPECTED_VERSION=$1
  if [[ $# > 1 ]]; then
    PLATFORM=$2
  fi
  if [[ $# > 2 ]]; then
    PI_TEST_SYSTEM=$3
  fi
fi

ROOT_DIR=$(realpath $(pwd)/$(dirname $0)/..)
SOLARCMD=${ROOT_DIR}/build/${PLATFORM}/solarcmd

loginfo "Checking against expected version $EXPECTED_VERSION"
ver=""
if [[ ${PLATFORM} == "linux_amd64" ]]; then
  ver=$(${SOLARCMD} -v -q 2>&1 | /bin/grep "version" | /bin/sed s/.*version.// | /bin/sed s/\".*//)
elif [[ ${PLATFORM} == "linux_arm" ]]; then
  scp ${SOLARCMD} sol@${PI_TEST_SYSTEM}:/tmp
  ver=$(ssh sol@${PI_TEST_SYSTEM} /tmp/solarcmd -v -q 2>&1 | /bin/grep "version" | /bin/sed s/.*version.// | /bin/sed s/\".*//)
else
  logerr "Unknown platform, $PLATFORM"
  exit 1
fi

if [[ $ver == $EXPECTED_VERSION ]]; then
  loginfo "Version matches ($ver)"
else 
  logerr "solarcmd version ($ver) does not match expected version ($EXPECTED_VERSION)"
  exit 1
fi
