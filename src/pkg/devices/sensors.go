package devices

// Sensor interface declares generic functions for a Sensor.
type Sensor interface {
	Device
  Resetter
}

// TemperatureSensor interface describes a temperature sensor.
type TemperatureSensor interface {
	Sensor
	Temperature() (float32, error)
}

// HumiditySensor interface describes a humidity sensor.
type HumiditySensor interface {
	Sensor
	Humidity() (float32, error)
}

// CurrentSensor interface describes a generic current sensor.
type CurrentSensor interface {
	Sensor
	Current() (float64, error)
}

// VoltageSensor interface describes a generic voltage sensor.
type VoltageSensor interface {
	Sensor
	Voltage() (float64, error)
}

// Shunt interface describes a device used to read the current through and the voltage across an external shunt.
type Shunt interface {
	CurrentSensor
  BusVoltage() (float64, error)
	ShuntVoltage() (float64, error)
}
