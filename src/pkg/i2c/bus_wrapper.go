package i2c

import (
	"fmt"
	"pkg/_vendor/i2c"
	"pkg/devices"
)

// MuxDefaultPortCount is the default number of ports of an I2C multiplexer.
const MuxDefaultPortCount = 8

// BusInterface abstracts the I2CBus (needs a new name).
type BusInterface interface {
	AddMultiplexer(addr devices.AddressType)
	RegisterSensor(s *Sensor) (err error)
	WriteDataWithReg(sid devices.DeviceID, reg, value byte) (err error)
	WriteData(sid devices.DeviceID, value byte) (err error)
	ReadByteBlock(sid devices.DeviceID, reg byte, readLength byte) (list []byte, err error)
	WriteByteBlock(sid devices.DeviceID, reg byte, list []byte) (err error)
}

// BusWrapper wraps a I2CBus to make it seem as if all devices are connected directly to the bus
// by automatically performing the switching through the multiplexers whenevrer data needs to be read or
// written to a device.
type BusWrapper struct {
	bus          *i2c.I2CBus
	multiplexers map[devices.AddressType]*Multiplexer // Map of Multiplexer by Multiplexer address
	sensors      map[devices.DeviceID]*Sensor
}

// Bus is a factory function used to create BusWrapper objects.
func Bus(busNum byte) (bus BusInterface, err error) {
	bw := &BusWrapper{}
	err = bw.Init(busNum, nil)
	if err != nil {
		return
	}
	bus = bw
	return
}

// Init function sets up a new BusWrapper. It needs a set of multiplexer addresses.
func (bw *BusWrapper) Init(busID byte, muxAddrList []devices.AddressType) (err error) {
	// create the underlying bus object
	bw.bus, err = i2c.Bus(busID)
	if err != nil {
		return
	}

	bw.multiplexers = make(map[devices.AddressType]*Multiplexer)
	bw.sensors = make(map[devices.DeviceID]*Sensor)

	for _, addr := range muxAddrList {
		bw.AddMultiplexer(addr)
	}
	return
}

// AddMultiplexer adds a multiplexer to the current I2C bus the given address.
func (bw *BusWrapper) AddMultiplexer(addr devices.AddressType) {
	mux := &Multiplexer{}
	mux.Init(addr, MuxDefaultPortCount, bw.bus)
	bw.multiplexers[addr] = mux
}

// RegisterSensor is used to register a sensor on the I2C bus if a sensor is connected directly to the bus.
func (bw *BusWrapper) RegisterSensor(s *Sensor) (err error) {
	bw.sensors[s.ID()] = s
	return
}

func (bw *BusWrapper) getSensorAddrFromID(sid devices.DeviceID) (addr devices.AddressType, err error) {
	s, ok := bw.sensors[sid]
	if !ok {
		err = fmt.Errorf("sensor with ID %d not found", sid)
	}
	addr = s.i2cAddress
	return
}

func (bw *BusWrapper) getSensor(sid devices.DeviceID) (s *Sensor, err error) {
	s, ok := bw.sensors[sid]
	if !ok {
		err = fmt.Errorf("Could not find sensor by ID, %d", sid)
		return nil, err
	}
	return s, nil
}

func (bw *BusWrapper) getMux(muxAddr devices.AddressType) (mux *Multiplexer, err error) {
	mux, ok := bw.multiplexers[muxAddr]
	if !ok {
		err = fmt.Errorf("Could not find mux by addr, %d", muxAddr)
		return nil, err
	}
	return mux, nil
}

func (bw *BusWrapper) switchMuxIfPresent(muxAddr devices.AddressType, muxPort uint) (err error) {
	// if muxAddr is non-zero, this means the sensor is connected to the bus via a mux
	if muxAddr != 0 {
		mux, err := bw.getMux(muxAddr)
		if err != nil {
			return err
		}
		err = mux.Switch(muxPort)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteDataWithReg writes a single byte to the given register at the given address.
func (bw *BusWrapper) WriteDataWithReg(sid devices.DeviceID, reg, value byte) (err error) {
	s, err := bw.getSensor(sid)
	if err != nil {
		return
	}

	err = bw.switchMuxIfPresent(s.muxAddr, s.muxPort)
	if err != nil {
		return err
	}

	return bw.bus.WriteByte(byte(s.i2cAddress), reg, value)
}

// WriteData writes a single byte to the given register at the given address.
func (bw *BusWrapper) WriteData(sid devices.DeviceID, value byte) (err error) {
	s, err := bw.getSensor(sid)
	if err != nil {
		return
	}

	err = bw.switchMuxIfPresent(s.muxAddr, s.muxPort)
	if err != nil {
		return err
	}

	return bw.bus.WriteByte2(byte(s.i2cAddress), value)
}

// ReadByteBlock reads a block of bytes from the given register at the given address.
func (bw *BusWrapper) ReadByteBlock(sid devices.DeviceID, reg byte, readLength byte) (list []byte, err error) {
	s, err := bw.getSensor(sid)
	if err != nil {
		return
	}

	err = bw.switchMuxIfPresent(s.muxAddr, s.muxPort)
	if err != nil {
		return nil, err
	}

	return bw.bus.ReadByteBlock(byte(s.i2cAddress), reg, readLength)
}

// WriteByteBlock writes a block of bytes to the given register at the given address.
func (bw *BusWrapper) WriteByteBlock(sid devices.DeviceID, reg byte, list []byte) (err error) {
	s, err := bw.getSensor(sid)
	if err != nil {
		return
	}

	err = bw.switchMuxIfPresent(s.muxAddr, s.muxPort)
	if err != nil {
		return err
	}

	return bw.bus.WriteByteBlock(byte(s.i2cAddress), reg, list)
}
