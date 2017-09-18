package models

import (
	"fmt"
	"time"

	"pkg/devices"
)

// HumidityMeasurement holds data read from a humidity sensor.
type HumidityMeasurement struct {
	ID               int       `db:"Id"`
	Timestamp        time.Time `db:"Timestamp"`
	HumiditySensorID int       `db:"HumiditySensorId"`
	RelativeHumidity float32   `db:"RelativeHumidity"`
}

func (m *HumidityMeasurement) String() (s string) {
	s += fmt.Sprintf("Timestamp: %s\n", m.Timestamp)
	s += fmt.Sprintf("RelativeHumidity: %.2f %%%%", m.RelativeHumidity)
	return s
}

// InsertIntoDb inserts the measurement into the database.
func (m *HumidityMeasurement) InsertIntoDb(dbmap devices.DbInserter) error {
	return dbmap.Insert(m)
}
