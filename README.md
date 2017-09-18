# Summary

## Build commands

* Build all (Runs **go build** for x64 arch)
  - (**F7**): `build/build_all.sh`
* Build for arm and upload to RPi (builds then copies both solarcmd exe and config files to pi-solar.local:bin/) 
  - `build/build_all.sh arm -f` 

## Test commands

* Unit tests with external server (Runs **goconvey test server**)
  - (**F8**): `build/run_test_server.sh`  
* Terminal task (test server):
  - (**Ctrl-Alt-T**)
* Manully run unit tests (Runs **go test** in each package directory that contains tests)
  - `build/run_unit_tests.sh`

## Debug commands
  
* Run solarcmd in debugger: 
  - (**F5**): Select "Launch Go" first


# Configure

## GOPATH

Make sure the path SolarMonitorController/src/ is appended to the GOPATH env. variable.

## Code structure

This is the structure of a go repository:

~~~~
$GOPATH/
  - src           <--- all source code goes here, each subdir for a separate package
    - cmd         <--- executable projects (source code)
    - pkg         <--- local packages (source code) 
      - logger
      - devices
      - db
    - 
  - pkg           <--- contains static libraries for built packages (used by other package and/or executable scripts)
  - bin           <--- contains executables (fully self-contained binaries, no external dependencies) 
~~~~

# Build, install, run go code

## Acquire (external) dependencies

~~~~
$ go get "github.com/go-gorp/gorp"
$ go get "github.com/Sirupsen/logrus"
$ go get "github.com/x-cray/logrus-prefixed-formatter"
$ go get "github.com/facebookgo/stack"
$ go get "github.com/spagettikod/gotracer"
$ go get "github.com/go-sql-driver/mysql"
~~~~

Optional packages:

~~~~
Test server: go get github.com/smartystreets/goconvey
JUnit generator (for Jenkins): go get -u github.com/jstemmer/go-junit-report
~~~~

## Common commands

* `./build/build_all.sh arm -f` (Build arm & upload to pi-solar.local:~/bin. Also copies solarcmd.config.json configuration file.) 
* `./build/build_all.sh` (Build amd64) 


## Go script can be run directly with the following command:

~~~~
$ go run src/scripts/solar_acquire_data.go
~~~~ 

## To build a go script into its own self-contained executable and install it into the bin directory:  

~~~~
$ cd src/scripts
$ go install 
~~~~

Binary will now appear in SolarMonitor/src/pi/scripts/solar/bin. It can be copied to another machine and run directly (without any go code, packages, etc).

## Cross-compile for arm

~~~~
$ GOOS=linux GOARCH=arm go build src/scripts/solar_acquire_data.go 
$ GOOS=linux GOARCH=arm go build src/scripts/solar_update_web_service.go 
~~~~

The cross-compiled arm binaries can be copied to the RPi.

