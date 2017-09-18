package i2c

import (
	"fmt"
	"pkg/_vendor/i2c"
	"pkg/devices"
	"pkg/logger"
)

// MuxID type - unqueue ID for a Multiplexer device.
type MuxID int

// PortID type - the ID of a port of a Multiplexer device (usually 0-7).
type PortID byte

var globalMuxID MuxID

// Multiplexer allows multiple I2C devices to be connected to the bus
type Multiplexer struct {
	id       MuxID
	address  devices.AddressType
	numPorts uint
	bus      *i2c.I2CBus

	Ports map[PortID]multiplexerPort
}

type multiplexerPort struct {
	PortID    PortID
	Connected bool
	SensorID  devices.DeviceID
}

// Init is used to initialise the Multiplexer object.
func (m *Multiplexer) Init(addr devices.AddressType, numPorts uint, b *i2c.I2CBus) {
	globalMuxID++
	m.id = globalMuxID
	m.address = addr
	m.numPorts = numPorts
	m.bus = b
}

// Switch is used to connect the a new device to the I2C bus (while automatically disconnecting the
// device that was previously connected).
func (m *Multiplexer) Switch(port uint) (err error) {
	if port >= m.numPorts {
		return fmt.Errorf("port num too large")
	}
	logger.Debugf("Switching to port: %d", port)
	return m.bus.WriteByte(byte(m.address), 0, (1 << port))
}
