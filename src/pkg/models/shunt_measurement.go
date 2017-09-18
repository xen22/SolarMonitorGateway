package models

import (
	"fmt"
	"pkg/devices"
	"time"
)

// ShuntMeasurement holds data related to a shunt in the circuit and which is read from an INA219 module.
type ShuntMeasurement struct {
	ID           int       `db:"Id"`
	CurrentAmps  float32   `db:"Current_A"`
	VoltageVolts float32   `db:"Voltage_V"`
	Timestamp    time.Time `db:"Timestamp"`
	ShuntID      int       `db:"ShuntId"`
	IntervalSecs int       `db:"Interval_s"`
}

func (m *ShuntMeasurement) String() (s string) {
	s += fmt.Sprintf("Timestamp: %s\n", m.Timestamp)
	s += fmt.Sprintf("Current: %.4f A\n", m.CurrentAmps)
	s += fmt.Sprintf("Voltage: %.4f V", m.VoltageVolts)
	return s
}

// InsertIntoDb inserts the measurement into the database.
func (m *ShuntMeasurement) InsertIntoDb(dbmap devices.DbInserter) error {
	return dbmap.Insert(m)
}
