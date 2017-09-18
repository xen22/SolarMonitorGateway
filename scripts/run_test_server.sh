#!/bin/bash
. $(dirname $0)/lib/logging

/home/ciprian/go/bin/goconvey --port 8082 --workDir /home/ciprian/dev/SolarMonitorController --excludedDirs solarcmd
