package devices

// DeviceFactory is used to instantiate different types of devices based on DeviceType and DeviceModel.
type DeviceFactory interface {
	Create(DeviceConfig) (Device, error)
}

// TemperatureSensorFactory is used to instantiate TemperatureSensors based on DeviceModel.
type TemperatureSensorFactory interface {
	Create(DeviceConfig) TemperatureSensor
}

// HumiditySensorFactory is used to instantiate HumiditySensors based on DeviceModel.
type HumiditySensorFactory interface {
	Create(DeviceConfig) HumiditySensor
}

// ShuntFactory is used to instantiate Shunts based on DeviceModel.
type ShuntFactory interface {
	Create(DeviceConfig) Shunt
}
