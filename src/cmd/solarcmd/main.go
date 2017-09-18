package main

import (
	"cmd/solarcmd/helpers"
	"flag"
	"os"
	"pkg/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// -------------------------------------------------------------------------------------------------------------------------
// Command line parameters configuration
// -------------------------------------------------------------------------------------------------------------------------

// Command line args that determine the mode of operation
// With the exception of debugMode, only one such command is executed (the rest are ignored, if mode than one is specified)

var (
	// command "modes"
	queryVersionMode          bool
	querySensorsMode          bool
	queryChargeControllerMode bool
	updateDbMode              bool
	updateWebServiceMode      bool

	// logging flags
	quietOutput  bool
	debugLogging bool

	// config file
	configFile string
)

func init() {
	// optional args
	flag.BoolVar(&quietOutput, "q", false, "quiet mode (suppress verbose logging info)")
	flag.BoolVar(&debugLogging, "d", false, "enable debug logging")
	flag.StringVar(&configFile, "c", "solarcmd.config.json", "specify a different configuration file (default is: solarcmd.config.json).")

	// mandatory args
	// Note: only a single one of the following should be specified at a time
	flag.BoolVar(&queryVersionMode, "v", false, "get current version")

	flag.BoolVar(&querySensorsMode, "S", false, "query all configured I2C sensors and retrieve data from them")
	flag.BoolVar(&queryChargeControllerMode, "C", false, "retrieve data from the charge controller")
	flag.BoolVar(&updateDbMode, "D", false, "save controller and sensor data to the database")
	flag.BoolVar(&updateWebServiceMode, "W", false, "push records from the db to the remote .Net web service")
}

// -------------------------------------------------------------------------------------------------------------------------

func queryVersion() {
	logger.Infof("%s version %s", os.Args[0], helpers.SolarCmdVersion)
}

func main() {
	logger.Infof("[%s] Begin.", time.Now())
	beginTime := time.Now()
	flag.Parse()

	if queryVersionMode {
		queryVersion()
		return
	}

	if querySensorsMode {
		logger.Info("Querying configured devices.")
		s, err := helpers.NewSystem(configFile)
		logger.FatalErrf(err, "Failed to create System object.")
		err = s.PrintMeasurements()
		logger.FatalErrf(err, "Failed to retrieve measurements.")
	} else if queryChargeControllerMode {

	} else if updateDbMode {
		s, err := helpers.NewSystem(configFile)
		logger.FatalErrf(err, "Failed to create System object.")
		err = s.SaveMeasurementsToDb()
		logger.FatalErrf(err, "Failed to save data to db.")
	} else if updateWebServiceMode {
		ws, err := helpers.NewWebService()
		logger.FatalErrf(err, "Failed to create WebService object.")
		err = ws.UpdateWebService()
		logger.FatalErrf(err, "Failed to create WebService object.")
	} else {
		flag.Usage()
		logger.Fatalf("No option mode specified. Exiting.")
	}

	logger.Infof("[%s] End. Total duration: %d sec.", time.Now(), time.Now().Unix()-beginTime.Unix())
}
