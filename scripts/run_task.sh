#!/bin/bash

export VSCODE_TASK=1
echo "[$(basename $0)]: sh -c \"$*\""
sh -c "$*"