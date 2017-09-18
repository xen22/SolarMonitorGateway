package adapters

import (
	"pkg/devices"
	"pkg/models"
)

// ShuntMeasurementFromShunt populates a ShuntMeasurement model from the Shunt.
func ShuntMeasurementFromShunt(s devices.Shunt) (sm models.ShuntMeasurement, err error) {
	current, err := s.Current()
	if err != nil {
		return
	}
	sm.CurrentAmps = float32(current)

	busVoltage, err := s.BusVoltage()
	if err != nil {
		return
	}
	shuntVoltage, err := s.BusVoltage()
	if err != nil {
		return
	}
	sm.VoltageVolts = float32(busVoltage + shuntVoltage)

	return
}
