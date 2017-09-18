package i2c

import (
	"pkg/devices"
	"pkg/logger"
)

// used to assign unique IDs to sensors
var idGenerator devices.DeviceID

func splitCmd(fullCmd uint16) (reg byte, cmd byte) {
	reg = byte(fullCmd >> 8)
	cmd = byte(fullCmd & 0xFF)
	return
}

// Sensor defines a generic sensor that is accessible via the I2C bus.
type Sensor struct {
	i2cAddress devices.AddressType
	muxAddr    devices.AddressType
	muxPort    uint
	bus        BusInterface
	devices.DeviceInfo
}

// Init is used to initialise an Sensor object.
func (s *Sensor) init(config devices.DeviceConfig) (err error) {
	err = s.DeviceInfo.Init(config)
	if err != nil {
		return
	}

	s.i2cAddress = config.Params.I2C.I2CAddress

	s.bus, err = Bus(config.Params.I2C.BusID)
	if err != nil {
		return
	}

	if config.Params.I2C.MuxAddress != 0 {
		s.muxAddr = config.Params.I2C.MuxAddress
		s.muxPort = config.Params.I2C.MuxPortNum
		s.bus.AddMultiplexer(s.muxAddr)
	}
	s.bus.RegisterSensor(s)
	return
}

func (s *Sensor) writeCommandWithReg(reg byte, cmd byte) (err error) {
	logger.Debugf("Cmd: 0x%X", cmd)
	return s.bus.WriteDataWithReg(s.ID(), reg, cmd)
}

func (s *Sensor) writeCommand(cmd byte) (err error) {
	logger.Debugf("Cmd: 0x%X", cmd)
	return s.bus.WriteData(s.ID(), cmd)
}
