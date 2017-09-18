#!/bin/bash
. $(dirname $0)/lib/logging

BUILD_CONFIG="Debug"
BUILD_VERSION="UnknownVersion"
PLATFORM="UnknownPlatform"

PROJECT_NAMES="solarcmd"
ROOT_DIR=$(realpath $(pwd)/$(dirname $0)/..)
OUTPUT_DIR="$ROOT_DIR/build"


if [[ $# -ne 3 ]]; then
  echo "Usage: $(basename $0) <BUILD_CONFIG> <BUILD_VERSION> <PLATFORM>"
  exit 1
else
  BUILD_CONFIG=$1
  BUILD_VERSION=$2
  PLATFORM=$3
fi

for proj in $PROJECT_NAMES ; do
  proj_dir=$(/usr/bin/find ${OUTPUT_DIR}/$PLATFORM/ -type f -name $proj)
  loginfo "Creating archive for project $proj found at \"$proj_dir\"."

  filename=${proj}-${BUILD_VERSION}-${BUILD_CONFIG}-${PLATFORM}
  new_dir=${OUTPUT_DIR}/${filename}
  mkdir ${new_dir}
  cp ${OUTPUT_DIR}/${PLATFORM}/${proj} ${new_dir} 
  cp ${ROOT_DIR}/src/cmd/${proj}/*.config.json ${new_dir}
  cd ${OUTPUT_DIR}
  tar zcf "${filename}.tar.gz" ${filename}
  cd -
  loginfo "Archive created: ${filename}."
  rm -rf ${new_dir}
done
