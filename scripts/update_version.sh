#!/bin/bash
. $(dirname $0)/lib/logging

NEW_VERSION="0.0.0.1"
if [[ $# -lt 1 ]]; then
  loginfo "Usage: $0 [<new_version>]"
  loginfo "Attempting to retrieve version from the most recent git tag on the current branch."
  NEW_VERSION=$(git tag -l --points-at HEAD | tail -1)
  if [[ $NEW_VERSION == "" ]]; then 
    logerr "Could not retrieve version."
    exit 1
  fi
else
  NEW_VERSION=$1
fi


ROOT_DIR=$(realpath $(pwd)/$(dirname $0)/..)
VERSION_FILE=${ROOT_DIR}/src/cmd/solarcmd/helpers/version.go

loginfo "Updating ${VERSION_FILE} file with new version ${NEW_VERSION}"
/bin/sed -i "/const SolarCmdVersion/c\const SolarCmdVersion = \"${NEW_VERSION}\"" ${VERSION_FILE}
