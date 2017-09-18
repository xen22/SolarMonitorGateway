package adapters

import (
	"encoding/json"
	"pkg/models"

	"github.com/spagettikod/gotracer"
)

// ChargeControllerMeasurementFromTracerStatus converts a gotracer.TracerStatus struct to ChargeControllerMeasurement
func ChargeControllerMeasurementFromTracerStatus(s gotracer.TracerStatus) (m models.ChargeControllerMeasurement, err error) {
	m = models.ChargeControllerMeasurement{}
	jsonData, err := json.Marshal(s)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		return
	}
	return m, nil
}
