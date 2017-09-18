package helpers

import (
	"cmd/solarcmd/config"
	"errors"
	"pkg/adapters"
	"pkg/db"
	"pkg/logger"
	"pkg/models"

	"time"

	"github.com/go-gorp/gorp"
)

func getChargeControllerID(dbMap *gorp.DbMap, siteName string) (id int64, err error) {
	query := "SELECT SolarMonitorDb.ChargeControllers.Id FROM SolarMonitorDb.ChargeControllers " +
		"INNER JOIN SolarMonitorDb.SolarSystems ON SolarSystemId = SolarMonitorDb.SolarSystems.Id " +
		"INNER JOIN SolarMonitorDb.Sites ON SiteId = SolarMonitorDb.Sites.Id " +
		"WHERE SolarMonitorDb.Sites.Name = \"" + siteName + "\" " +
		"LIMIT 1"
	id, err = dbMap.SelectInt(query)
	return
}

func getShuntID(dbMap *gorp.DbMap, siteName string, shuntName string) (id int64, err error) {
	query := "SELECT SolarMonitorDb.Shunts.Id FROM SolarMonitorDb.Shunts " +
		"INNER JOIN SolarMonitorDb.SolarSystems ON SolarSystemId = SolarMonitorDb.SolarSystems.Id " +
		"INNER JOIN SolarMonitorDb.Sites ON SiteId = SolarMonitorDb.Sites.Id " +
		"WHERE SolarMonitorDb.Sites.Name = \"" + siteName + "\" " +
		"AND SolarMonitorDb.Shunts.Name = \"" + shuntName + "\" "
	id, err = dbMap.SelectInt(query)
	if id == 0 {
		err = errors.New("Failed to find Shunt in database SolarMonitorDb.")
		return
	}
	return
}

func getTemperatureSensorID(dbMap *gorp.DbMap, siteName string, tempSensorName string) (id int64, err error) {
	query := "SELECT SolarMonitorDb.TemperatureSensors.Id FROM SolarMonitorDb.TemperatureSensors " +
		"INNER JOIN SolarMonitorDb.WeatherBases ON WeatherBaseId = SolarMonitorDb.WeatherBases.Id " +
		"INNER JOIN SolarMonitorDb.Sites ON SiteId = SolarMonitorDb.Sites.Id " +
		"WHERE SolarMonitorDb.Sites.Name = \"" + siteName + "\" " +
		" AND SolarMonitorDb.TemperatureSensors.Name = \"" + tempSensorName + "\" "
	id, err = dbMap.SelectInt(query)
	if id == 0 {
		err = errors.New("Failed to find TemperatureSensor in database SolarMonitorDb.")
		return
	}
	return
}

func getHumiditySensorID(dbMap *gorp.DbMap, siteName string, humSensorName string) (id int64, err error) {
	query := "SELECT SolarMonitorDb.HumiditySensors.Id FROM SolarMonitorDb.HumiditySensors " +
		"INNER JOIN SolarMonitorDb.WeatherBases ON WeatherBaseId = SolarMonitorDb.WeatherBases.Id " +
		"INNER JOIN SolarMonitorDb.Sites ON SiteId = SolarMonitorDb.Sites.Id " +
		"WHERE SolarMonitorDb.Sites.Name = " + siteName + " " +
		" AND SolarMonitorDb.HumiditySensors.Name = \"" + humSensorName + "\" "
	id, err = dbMap.SelectInt(query)
	if id == 0 {
		err = errors.New("Failed to find HumiditySensor in database SolarMonitorDb.")
		return
	}
	return
}

// UpdateDatabase creates and insers into the database a SolarRecord
func UpdateDatabase() {

	logger.Info("Get data from charge controller.")
	status, err := QueryChargeController()
	logger.FatalErrf(err, "Could not get status data from the charge controller.")

	logger.Info("Creating sensor objects.")
	tSensors, hSensors, shunts, err := CreateSensors()
	logger.FatalErrf(err, "Failed to create sensors")

	dbmap := db.InitDb(config.DbName, config.DbUser, config.DbPassword)
	defer func() {
		logger.Info("Closing db connection.")
		dbmap.Db.Close()
	}()

	// Configure tables
	dbmap.AddTableWithName(adapters.TracerStatusDbAdapter{}, config.ChargeControllerMeasurementsTable).SetKeys(false, "Id")
	dbmap.AddTableWithName(models.ShuntMeasurement{}, config.ShuntMeasurementsTable).SetKeys(false, "Id")
	dbmap.AddTableWithName(models.TemperatureMeasurement{}, config.TemperatureMeasurementsTable).SetKeys(false, "Id")
	dbmap.AddTableWithName(models.HumidityMeasurement{}, config.HumidityMeasurementsTable).SetKeys(false, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	logger.FatalErrf(err, "Create tables failed")

	dbstatus := adapters.DbAdapterfromTracerStatus(status)
	ccid, err := getChargeControllerID(dbmap, config.SiteName)
	logger.FatalErrf(err, "Unable to retrieve ChargeController ID from db.")
	dbstatus.ChargeControllerID = int(ccid)

	logger.Info("Saving charge controller data to db.")
	err = dbmap.Insert(&dbstatus)
	logger.FatalErrf(err, "Unable to insert row in table ChargeControllerMeasurements. Err: %s.", err)

	logger.Info("Reading data from shunt sensor.")
	for shuntName, s := range shunts {
		sm, err := adapters.ShuntMeasurementFromShunt(s)
		logger.FatalErrf(err, "Unable to convert read data from shunt sensor %s.", shuntName)

		shuntID, err := getShuntID(dbmap, config.SiteName, shuntName)
		logger.FatalErrf(err, "Unable to obtain Shunt ID from the database for shunt sensor %s.", shuntName)
		sm.Timestamp = time.Now()
		logger.Debugf("shunt id: %d", shuntID)
		sm.ShuntID = int(shuntID)
		logger.Debugf("ShuntMeasurement: %s", sm)

		logger.Infof("Saving data for shunt %q to db.", shuntName)
		err = dbmap.Insert(&sm)
		logger.FatalErrf(err, "Unable to insert row. ")
	}

	logger.Info("Reading data from temperature sensor.")
	for tName, t := range tSensors {
		tm, err := adapters.TemperatureMeasurementFromSensor(t)
		logger.FatalErrf(err, "Unable to convert read data from temperature sensor %s.", tName)

		tempSensorID, err := getTemperatureSensorID(dbmap, config.SiteName, tName)
		logger.FatalErrf(err, "Unable to obtain TemperatureSensor ID from the database for temperature sensor %s.", tName)
		tm.Timestamp = time.Now()
		tm.TemperatureSensorID = int(tempSensorID)

		logger.Info("Saving temperature sensor data to db.")
		err = dbmap.Insert(&tm)
		logger.FatalErrf(err, "Unable to insert row. ")
	}

	logger.Info("Reading data from humidity sensor.")
	for hName, h := range hSensors {
		hm, err := adapters.HumidityMeasurementFromSensor(h)
		logger.FatalErrf(err, "Unable to convert read data from humidity sensor %s.", hName)

		humSensorID, err := getHumiditySensorID(dbmap, config.SiteName, hName)
		logger.FatalErrf(err, "Unable to obtain HumiditySensor ID from the database for humidity sensor %s.", hName)
		hm.Timestamp = time.Now()
		hm.HumiditySensorID = int(humSensorID)

		logger.Info("Saving humidity sensor data to db.")
		err = dbmap.Insert(&hm)
		logger.FatalErrf(err, "Unable to insert row. ")
	}
}
