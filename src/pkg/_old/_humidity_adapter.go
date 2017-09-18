package adapters

import (
	"pkg/devices"
	"pkg/models"
)

// HumidityMeasurementFromSensor populates a HumidityMeasurement model from the corresponding i2c sensor object.
func HumidityMeasurementFromSensor(h devices.HumiditySensor) (hm models.HumidityMeasurement, err error) {
	hum, err := h.Humidity()
	if err != nil {
		return
	}
	hm.RelativeHumidity = float32(hum)

	return
}
