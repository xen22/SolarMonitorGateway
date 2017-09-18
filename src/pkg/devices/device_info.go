package devices

// DeviceInfo contains data describing a Device
type DeviceInfo struct {
	name    string
	id      DeviceID
	devType SensorType
	//	SensorData
}

// Init is used to initialise a DeviceInfo object.
func (d *DeviceInfo) Init(config DeviceConfig) (err error) {
	d.name = config.Name
	d.id = config.DbID
	d.devType = config.Type
	return nil
}

// Name returns the name of the device.
func (d *DeviceInfo) Name() string {
	return d.name
}

// ID returns the device ID.
func (d *DeviceInfo) ID() DeviceID {
	return d.id
}

// Type returns the type of the device.
func (d *DeviceInfo) Type() SensorType {
	return d.devType
}
