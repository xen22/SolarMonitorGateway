// Converted (with modifications) from the python implementation of Adafruit_SHT31.py
// Source:https://github.com/ralf1070/Adafruit_Python_SHT31
// Original authors and licence below:

// # Copyright (c) 2014 Adafruit Industries
// # Author: Tony DiCola
// #
// # Based on the BME280 driver with SHT31D changes provided by
// # Ralf Mueller, Erfurt
// #
// # Permission is hereby granted, free of charge, to any person obtaining a copy
// # of this software and associated documentation files (the "Software"), to deal
// # in the Software without restriction, including without limitation the rights
// # to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// # copies of the Software, and to permit persons to whom the Software is
// # furnished to do so, subject to the following conditions:
// #
// # The above copyright notice and this permission notice shall be included in
// # all copies or substantial portions of the Software.
// #
// # THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// # IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// # FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// # AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// # LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// # OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// # THE SOFTWARE.

package i2c

import (
	"errors"
	"fmt"
	"pkg/devices"
	"time"
)

const (
	softReset      = 0x30A2
	measureHighRep = 0x2400
)

// SHT31 struct allows access to an Adafruit SHT31-D temperature and humidity sensor.
type SHT31 struct {
	Sensor
	//DeviceInfo
}

// NewSHT31Sensor returns a new instance of SHT31 type sensor.
func NewSHT31Sensor(config devices.DeviceConfig) (s *SHT31, err error) {
	s = &SHT31{}
	err = s.init(config)
	if err != nil {
		return nil, err
	}
	return
}

// Reset performs a soft reset of the device.
func (s *SHT31) Reset() (err error) {
	return s.writeCommandWithReg(splitCmd(softReset))
}

func (s *SHT31) crc8(buffer []byte) (err byte) {
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

func (s *SHT31) readTempAndHumidity() (temp float32, hum float32, err error) {
	err = s.writeCommandWithReg(splitCmd(measureHighRep))
	if err != nil {
		return
	}

	time.Sleep(15 * time.Millisecond)
	buffer, err := s.bus.ReadByteBlock(s.ID(), 0, 6)
	if err != nil {
		return
	}

	if buffer[2] != s.crc8(buffer[0:2]) {
		err = errors.New("IO error: crc8  check failed")
		return
	}

	rawTemp := float32(int16(buffer[0])<<8 | int16(buffer[1]))
	temp = 175.0*rawTemp/0xFFFF - 45.0

	if buffer[5] != s.crc8(buffer[3:5]) {
		err = errors.New("IO error: crc8  check failed")
		return
	}

	rawHumidity := float32(int16(buffer[3])<<8 | int16(buffer[4]))
	hum = 100.0*rawHumidity/0xFFFF + 100

	return
}

// Temperature is used to obtain the current temperature reading from the sensor.
func (s *SHT31) Temperature() (temp float32, err error) {
	temp, _, err = s.readTempAndHumidity()
	return
}

// Humidity is used to obtain the current humidity reading from the sensor.
func (s *SHT31) Humidity() (h float32, err error) {
	_, h, err = s.readTempAndHumidity()
	return
}

// NewMeasurement returns a new set of data from the sensor representing a Measurement unit.
func (s *SHT31) NewMeasurement() (m devices.Measurement, err error) {
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
