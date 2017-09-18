package devices

import "fmt"

// SensorType describes the type of sensor.
type SensorType int

const (
	// TemperatureType is a type identifier for a temperature sensor.
	TemperatureType SensorType = iota

	// HumidityType is a type identifier for a humidity sensor.
	HumidityType SensorType = iota

	// CurrentType is a type identifier for a current sensor.
	CurrentType SensorType = iota

	// ShuntType is a type identifier for a shunt sensor.
	ShuntType SensorType = iota

	// ChargeControllerType is a type identifier for a charge controller.
	ChargeControllerType SensorType = iota
)

// UnmarshalJSON converts string to SensorType.
func (s *SensorType) UnmarshalJSON(value []byte) (err error) {
	switch string(value) {
	case "TemperatureType":
		*s = TemperatureType
	case "HumidityType":
		*s = HumidityType
	case "CurrentType":
		*s = CurrentType
	case "ShuntType":
		*s = ShuntType
	case "ChargeControllerType":
		*s = ChargeControllerType
	default:
		err = fmt.Errorf("Unknown sensor type %s", value)
	}
	return
}
