package devices

// DeviceID type
type DeviceID int

// DeviceConfig defines the configuration parameters for a device.
type DeviceConfig struct {
	Type        SensorType `json:"Type,string"`
	Model       ModelType  `json:"Model,string"`
	Name        string
	Description string
	Location    string
	DbID        DeviceID `json:"DbID,int"` // This must match the ID of the corresponding device/sensor in the database
	Params      DeviceSpecificParams
}

// DeviceSpecificParams is a struct containing any device specific parameters.
type DeviceSpecificParams struct {
	I2C    I2CParams
	INA219 INA219Params
	Tracer TracerParams
}

// I2CParams contains parameters specific to I2C sensors.
type I2CParams struct {
	BusID      byte
	I2CAddress AddressType `json:"I2CAddress,string"`
	MuxAddress AddressType `json:"MuxAddress,string"`
	MuxPortNum uint
}

// INA219Params contains parameters specific to a INA219 current sensor with an external shunt.
type INA219Params struct {
	ShuntResistorMilliOhms     float64 // size of the external shunt resistor (assuming the 100 mOhm INA219 resistor has been removed)
	MaxCurrentA                int     // max expected current through the circuit to be measured (not the max current spec of the shunt)
	MaxShuntVoltageV           float64
	MaxBusVoltageV             float64
	CurrentResolutionMilliAmps float64 // this is the configured curr. resolution used to calculate the calibration value
	// note though that the INA219 will calculate the actual current resolution as "10uV / shunt resistor"
}

// TracerParams contains parameters specific to a Tracer charge controller.
type TracerParams struct {
	SerialDevice string
}
