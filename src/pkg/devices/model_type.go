package devices

import "fmt"

// ModelType describes the model of a sensor.
type ModelType int

const (
	// SHT31Model refers to the SHT31 temperature and humidity sensor.
	SHT31Model ModelType = iota
	// HTU21Model refers to the SHT31 temperature and humidity sensor.
	HTU21Model ModelType = iota
	// INA219Model refers to the SHT31 temperature and humidity sensor.
	INA219Model ModelType = iota
	// Tracer4215BNModel refers to the Tracer4215BNModel charge controller.
	Tracer4215BNModel ModelType = iota
	// SimulatedTracer4215BNModel refers to the dummy SimulatedTracer4215BNModel change controller.
	SimulatedTracer4215BNModel ModelType = iota
)

// UnmarshalJSON converts string to ModelType.
func (m *ModelType) UnmarshalJSON(value []byte) (err error) {
	switch string(value) {
	case "SHT31":
		*m = SHT31Model
	case "HTU21":
		*m = HTU21Model
	case "INA219":
		*m = INA219Model
	case "Tracer4215BN":
		*m = Tracer4215BNModel
	case "SimulatedTracer4215BN":
		*m = SimulatedTracer4215BNModel
	default:
		err = fmt.Errorf("Unknown model type %s", value)
	}
	return
}
