package i2c

import (
	"fmt"
	"pkg/devices"
	"pkg/logger"
	"pkg/models"
	"time"
)

const (
	// HTU21D-F Commands
	readTempCmd     = 0xE3
	readHumidityCmd = 0xE5
	//writeRegisterCmd = 0xE6
	readRegisterCmd = 0xE7
	resetCmd        = 0xFE
	resetResp       = 0x02

	readRegister = 0xE7
)

// HTU21 struct allows access to an Adafruit HTU21D-F temperature and humidity sensor.
type HTU21 struct {
	Sensor
}

// NewHTU21Sensor returns a new instance of HTU21 type sensor.
func NewHTU21Sensor(config devices.DeviceConfig) (s *HTU21, err error) {
	s = &HTU21{}
	err = s.init(config)
	if err != nil {
		return nil, err
	}
	return
}

// Reset performs a soft reset of the device.
func (s *HTU21) Reset() (err error) {
	err = s.writeCommand(resetCmd) // 0xE6FE
	if err != nil {
		return
	}

	logger.Debug("Sending command readRegisterCmd")
	err = s.writeCommand(readRegisterCmd) // 0xE6E7
	if err != nil {
		return
	}

	time.Sleep(55 * time.Millisecond)

	logger.Debug("Reading byte from reg")
	var ret []byte
	ret, err = s.bus.ReadByteBlock(s.ID(), readRegister, 1) // reg: E7
	if err != nil {
		return
	}
	if ret[0] != resetResp {
		logger.Debugf("ret: %v", ret)
		err = fmt.Errorf("Return byte after reset incorrect, expected 0x%X, got 0x%X", resetResp, ret[0])
	}

	return
}

func (s *HTU21) crc8(buffer []byte) (err byte) {

	polynomial := byte(0x31)
	crc := byte(0xFF)

	for i := 0; i < len(buffer); i++ {
		crc ^= buffer[i]
		for j := 8; j > 0; j-- {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc = (crc << 1)
			}
		}
	}
	return crc & 0xFF
}

// Temperature is used to obtain the current temperature reading from the sensor.
func (s *HTU21) Temperature() (temp float32, err error) {

	err = s.writeCommand(readTempCmd)
	if err != nil {
		return
	}
	time.Sleep(55 * time.Millisecond) // > 50ms

	buffer, err := s.bus.ReadByteBlock(s.ID(), readTempCmd, 3)
	if err != nil {
		return
	}
	logger.Debugf("Raw bytes for temp: %v", buffer)

	rawTemp := float32(uint16(buffer[0])<<8 | uint16(buffer[1]))
	temp = rawTemp/65536*175.72 - 46.85
	return
}

// Humidity is used to obtain the current humidity reading from the sensor.
func (s *HTU21) Humidity() (hum float32, err error) {
	err = s.writeCommand(readHumidityCmd)
	if err != nil {
		return
	}
	time.Sleep(55 * time.Millisecond) // > 50ms

	buffer, err := s.bus.ReadByteBlock(s.ID(), readHumidityCmd, 3)
	if err != nil {
		return
	}
	logger.Debugf("Raw bytes for hum: %v", buffer)

	rawHum := float32(uint16(buffer[0])<<8 | uint16(buffer[1]))
	hum = rawHum/65536*125 - 6
	return
}

func newTemperatureMeasurement(s devices.TemperatureSensor, sensorID int) (m devices.Measurement, err error) {
	tm := &models.TemperatureMeasurement{}
	temperature, err := s.Temperature()
	if err != nil {
		return nil, err
	}
	tm.TemperatureCelsius = temperature
	tm.TemperatureSensorID = int(sensorID)
	tm.Timestamp = time.Now()
	m = tm

	return m, nil
}

func newHumidityMeasurement(s devices.HumiditySensor, sensorID int) (m devices.Measurement, err error) {
	hm := &models.HumidityMeasurement{}
	humidity, err := s.Humidity()
	if err != nil {
		return nil, err
	}
	hm.RelativeHumidity = humidity
	hm.HumiditySensorID = int(sensorID)
	hm.Timestamp = time.Now()

	m = hm
	return m, nil
}

// NewMeasurement returns a new set of data from the sensor representing a Measurement unit.
func (s *HTU21) NewMeasurement() (m devices.Measurement, err error) {
	if s.Type() == devices.TemperatureType {
		m, err = newTemperatureMeasurement(s, int(s.ID()))
	} else if s.Type() == devices.HumidityType {
		m, err = newHumidityMeasurement(s, int(s.ID()))
	} else {
		err = fmt.Errorf("Unknown sensor type %d", s.Type())
		return
	}
	return m, nil
}
