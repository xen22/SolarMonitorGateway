package devices

// Device is a generic device in the system, such as an I2C sensor, ChargeController, etc.
type Device interface {
	Name() string
	ID() DeviceID
	//Description() string

	NewMeasurement() (Measurement, error)
}
