package db

import (
	"cmd/solarcmd/config"
	"database/sql"
	"fmt"
	"pkg/devices"
	"pkg/logger"
	"pkg/models"

	"github.com/go-gorp/gorp"
)

// SolarDb allows inserting and retrieving Measurements and other data from the db.
type SolarDb struct {
	dbConfig devices.DatabaseConfig
	dbMap    *gorp.DbMap
}

// NewSolarDb is used to create a new SolarDb instance.
func NewSolarDb(dbConfig devices.DatabaseConfig) (sdb *SolarDb) {
	sdb = &SolarDb{}
	sdb.dbConfig = dbConfig
	return sdb
}

// Init opens a connection to a database, creates and configures a gorp.DbMap and returns it.
func (sdb *SolarDb) Init() error {
	logger.Infof("Open db connection (db: %q)", sdb.dbConfig.Name)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", sdb.dbConfig.User, sdb.dbConfig.Password, sdb.dbConfig.Name))
	if err != nil {
		return err
	}

	// construct a gorp DbMap
	sdb.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}

	sdb.dbMap.AddTableWithName(models.ChargeControllerMeasurement{}, config.ChargeControllerMeasurementsTable).SetKeys(false, "Id")
	sdb.dbMap.AddTableWithName(models.ShuntMeasurement{}, config.ShuntMeasurementsTable).SetKeys(false, "Id")
	sdb.dbMap.AddTableWithName(models.TemperatureMeasurement{}, config.TemperatureMeasurementsTable).SetKeys(false, "Id")
	sdb.dbMap.AddTableWithName(models.HumidityMeasurement{}, config.HumidityMeasurementsTable).SetKeys(false, "Id")

	return nil
}

// Close closes the database.
func (sdb *SolarDb) Close() error {
	if sdb.dbMap != nil {
		return sdb.dbMap.Db.Close()
	}
	return nil
}

// InsertMeasurement is used to insert any type of Measurement into the database.
func (sdb *SolarDb) InsertMeasurement(m devices.Measurement) error {
	sdb.dbMap.Insert(m)
	return nil
}
