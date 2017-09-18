package models

import (
	"fmt"
	"pkg/devices"
	"time"
)

// TemperatureMeasurement holds data read from a temperature sensor.
type TemperatureMeasurement struct {
	ID                  int       `db:"Id"`
	Timestamp           time.Time `db:"Timestamp"`
	TemperatureSensorID int       `db:"TemperatureSensorId"`
	TemperatureCelsius  float32   `db:"Temperature_C"`
}

func (m *TemperatureMeasurement) String() (s string) {
	s += fmt.Sprintf("Timestamp: %s\n", m.Timestamp)
	s += fmt.Sprintf("Temperature: %.2f C", m.TemperatureCelsius)
	return s
}

// InsertIntoDb inserts the measurement into the database.
func (m *TemperatureMeasurement) InsertIntoDb(dbmap devices.DbInserter) error {
	return dbmap.Insert(m)
}
