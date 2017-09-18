package factories

import (
	"fmt"
	"pkg/devices"
	"pkg/i2c"
	"pkg/sim"
	"pkg/solar"
)

// DeviceFactory is device factory implementation.
type DeviceFactory struct {
}

// Create creates a new device based on type and model.
func (f *DeviceFactory) Create(c devices.DeviceConfig) (d devices.Device, err error) {
	switch c.Type {
	case devices.ShuntType:
		switch c.Model {
		case devices.INA219Model:
			d, err = i2c.NewINA219Sensor(c)
		default:
			err = fmt.Errorf("Unknown device model %d for type %d", c.Model, c.Type)
		}

	case devices.TemperatureType:
		fallthrough
	case devices.HumidityType:
		switch c.Model {
		case devices.HTU21Model:
			d, err = i2c.NewHTU21Sensor(c)
		case devices.SHT31Model:
			d, err = i2c.NewSHT31Sensor(c)
		default:
			err = fmt.Errorf("Unknown device model %d for type %d", c.Model, c.Type)
		}

	case devices.ChargeControllerType:
		switch c.Model {
		case devices.Tracer4215BNModel:
			d, err = solar.NewTracer4215BN(c)
		case devices.SimulatedTracer4215BNModel:
			d, err = sim.NewSimulatedTracer4215BN(c)
		default:
			err = fmt.Errorf("Unknown device model %d for type %d", c.Model, c.Type)
		}

	default:
		err = fmt.Errorf("Unknown device type %d", c.Type)
	}
	return
}

// NewTemperatureSensor is a factory that creates sensors and returns them via the TemperatureSensor interface.
func NewTemperatureSensor(config devices.DeviceConfig) (s devices.TemperatureSensor, err error) {
	if config.Type != devices.TemperatureType {
		err = fmt.Errorf("Invalid sensor type, %d", config.Type)
		return
	}

	switch config.Model {
	case devices.SHT31Model:
		s, err = i2c.NewSHT31Sensor(config)

	case devices.HTU21Model:
		s, err = i2c.NewHTU21Sensor(config)

	default:
		err = fmt.Errorf("Unknown sensor model, %d", config.Model)
	}
	return
}

// NewHumiditySensor is a factory that creates sensors and returns them via the HumiditySensor interface.
func NewHumiditySensor(config devices.DeviceConfig) (s devices.HumiditySensor, err error) {
	if config.Type != devices.HumidityType {
		err = fmt.Errorf("Invalid sensor type, %d", config.Type)
		return
	}

	switch config.Model {
	case devices.SHT31Model:
		s, err = i2c.NewSHT31Sensor(config)

	case devices.HTU21Model:
		s, err = i2c.NewHTU21Sensor(config)

	default:
		err = fmt.Errorf("Unknown sensor model, %d", config.Model)
	}
	return
}

// NewShuntSensor is a factory that creates sensors and returns them via the Shunt interface.
func NewShuntSensor(config devices.DeviceConfig) (s devices.Shunt, err error) {
	if config.Type != devices.ShuntType {
		err = fmt.Errorf("Invalid sensor type, %d", config.Type)
		return
	}

	switch config.Model {
	case devices.INA219Model:
		var shunt *i2c.INA219
		shunt, err = i2c.NewINA219Sensor(config)
		if err != nil {
			return
		}
		s = shunt

	default:
		err = fmt.Errorf("Unknown sensor model, %d", config.Model)
	}
	return
}
