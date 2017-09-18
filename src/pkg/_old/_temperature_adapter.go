package adapters

import (
	"pkg/devices"
	"pkg/models"
)

// TemperatureMeasurementFromSensor populates a TemperatureMeasurement model from the corresponding i2c sensor object.
func TemperatureMeasurementFromSensor(t devices.TemperatureSensor) (tm models.TemperatureMeasurement, err error) {
	temp, err := t.Temperature()
	if err != nil {
		return
	}
	tm.TemperatureCelsius = float32(temp)

	return
}
