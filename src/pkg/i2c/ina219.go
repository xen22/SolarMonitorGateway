// Adapted from Adafruit_INA219/Adafruit_INA219.h
// Source: https://github.com/adafruit/Adafruit_INA219/blob/master/Adafruit_INA219.h
// Original author   K. Townsend (Adafruit Industries), License  BSD

package i2c

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"pkg/devices"
	"pkg/logger"
	"pkg/models"
	"time"
)

const (
	i2cAddress01 = 0x40
	i2cAddress02 = 0x41

	configResetCmd = 0x8000 // Reset Bit

	configBusVoltageRangeMask = 0x2000 // Bus Voltage Range Mask
	configBusVoltageRange16V  = 0x0000 // 0-16V Range
	configBusVoltageRange32V  = 0x2000 // 0-32V Range

	configGainMask    = 0x1800 // Gain Mask
	configGain1_40mV  = 0x0000 // Gain 1, 40mV Range
	configGain2_80mV  = 0x0800 // Gain 2, 80mV Range
	configGain4_160mV = 0x1000 // Gain 4, 160mV Range
	configGain8_320mV = 0x1800 // Gain 8, 320mV Range

	configBusADCResMask  = 0x0780 // Bus ADC Resolution Mask
	configBusADCRes9Bit  = 0x0080 // 9-bit bus res = 0..511
	configBusADCRes10Bit = 0x0100 // 10-bit bus res = 0..1023
	configBusADCRes11Bit = 0x0200 // 11-bit bus res = 0..2047
	configBusADCRes12Bit = 0x0400 // 12-bit bus res = 0..4097

	configShuntADCResMask                = 0x0078 // Shunt ADC Resolution and Averaging Mask
	configShuntADCRes9Bit1Sample84us     = 0x0000 // 1 x 9-bit shunt sample
	configShuntADCRes10Bit1Sample148us   = 0x0008 // 1 x 10-bit shunt sample
	configShuntADCRes11Bit1Sample276us   = 0x0010 // 1 x 11-bit shunt sample
	configShuntADCRes12Bit1Sample532us   = 0x0018 // 1 x 12-bit shunt sample
	configShuntADCRes12Bit2Sample1060us  = 0x0048 // 2 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit4Sample2130us  = 0x0050 // 4 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit8Sample4260us  = 0x0058 // 8 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit16Sample8510us = 0x0060 // 16 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit32Sample17ms   = 0x0068 // 32 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit64Sample34ms   = 0x0070 // 64 x 12-bit shunt samples averaged together
	configShuntADCRes12Bit128Sample69ms  = 0x0078 // 128 x 12-bit shunt samples averaged together

	configModeMask                         = 0x0007 // Operating Mode Mask
	configModePowerDown                    = 0x0000
	configModeShuntVoltageTriggered        = 0x0001
	configModeBusVoltageTriggered          = 0x0002
	configModeShuntAndBusVoltageTriggered  = 0x0003
	configModeADCOff                       = 0x0004
	configModeShuntVoltageContinuous       = 0x0005
	configModeBusVoltageContinuous         = 0x0006
	configModeShuntAndBusVoltageContinuous = 0x0007

	regConfig       = 0x00
	regShuntVoltage = 0x01
	regBusVoltage   = 0x02
	regPower        = 0x03
	regCurrent      = 0x04
	regCalibration  = 0x05

	// INA219 has the following fixed resolutions (LSBs)
	busVoltageLSBmV   = 4.0
	shuntVoltageLSBuV = 10.0
)

// INA219 implements functionality to allow reading voltage and current from a INA219 sensor.
// Shunt specific configuration information needs to be passed in via SensorSpecificParams (see configureSensor())
// Notes:
// * Current LSB (resolution) calculation method:
// 1. Calculate possible range of LSBs (Min = 15-bit, Max = 12-bit)
// This example is for a 20A current
// MinimumLSB = MaxExpected_I/32767 = 20/32767
// MinimumLSB = 0.0006104 (0.61 mA per bit) ==> approximate up to 1 mA
// MaximumLSB = MaxExpected_I/4096
// MaximumLSB = 0.004882 (4.882mA per bit)
// 2. Choose an LSB between the min and max values
// 	(Preferrably a roundish number close to MinLSB)
// CurrentLSB = 1.0    (1mA per bit)
type INA219 struct {
	Sensor
	INA219Params

	CalibrationValue uint16
	Configuration    uint16
	isInitialised    bool
	currentAdj       float64
	busVoltageAdj    float64
	shuntVoltageAdj  float64
}

// INA219Params contains parameters specific to the INA219 sensor
type INA219Params struct {
	ShuntResistorMilliOhms     float64 // size of the external shunt resistor (assuming the 100 mOhm INA219 resistor has been removed)
	MaxCurrentA                int     // max expected current through the circuit to be measured (not the max current spec of the shunt)
	MaxShuntVoltageV           float64
	MaxBusVoltageV             float64
	CurrentResolutionMilliAmps float64 // this is the configured curr. resolution used to calculate the calibration value
	// note though that the INA219 will calculate the actual current resolution as "10uV / shunt resistor"
}

// NewINA219Sensor is a factory function.
func NewINA219Sensor(config devices.DeviceConfig) (s *INA219, err error) {
	s = &INA219{}

	err = s.init(config)
	if err != nil {
		return nil, err
	}

	s.configureSensor(config.Params)
	return
}

func (s *INA219) printConfig() {
	logger.Debugf("== INA219 sensor: %s ==============", s.Name())
	logger.Debugf("-- MaxExpectedCurrent:          %d A", s.MaxCurrentA)
	logger.Debugf("-- MaxBusVoltageV:              %.2f V", s.MaxBusVoltageV)
	logger.Debugf("-- MaxShuntVoltageV:            %.2f V", s.MaxShuntVoltageV)
	logger.Debugf("-- ShuntResistorMilliOhms:      %.2f mOhm", s.ShuntResistorMilliOhms)
	logger.Debugf("-- CurrentResolutionMilliAmps:  %.2f mA", s.CurrentResolutionMilliAmps)
	logger.Debugf("-- ShuntResistorMilliOhms:      %.2f mOhm", s.ShuntResistorMilliOhms)

	logger.Debugf("-- Configuration:               %b", s.Configuration)
	logger.Debugf("-- CalibrationValue:            0x%X", s.CalibrationValue)
	logger.Debugf("-- currentAdj:                  %.4f", s.currentAdj)
	logger.Debugf("-- busVoltageAdj:               %.4f", s.busVoltageAdj)
	logger.Debugf("-- shuntVoltageAdj:             %.4f", s.shuntVoltageAdj)

	logger.Debugf("-- Bus voltage resolution:      %.2f mV", busVoltageLSBmV)
	logger.Debugf("-- Shunt voltage resolution:    %.2f uV", shuntVoltageLSBuV)
	logger.Debugf("-- Current resolution (actual): %.2f mA", shuntVoltageLSBuV/s.ShuntResistorMilliOhms)
}

func (s *INA219) calculateCalibrationValue() (cal uint16, err error) {
	val := uint64(math.Trunc(0.04096 / (s.CurrentResolutionMilliAmps / 1000 * s.ShuntResistorMilliOhms / 1000)))
	if val > math.MaxUint16 {
		err = fmt.Errorf("Max value of the calibration register exceeded, calculated value: %X", val)
		return
	}
	return uint16(val), nil
}

func (s *INA219) configureSensor(params devices.DeviceSpecificParams) (err error) {

	s.ShuntResistorMilliOhms = params.INA219.ShuntResistorMilliOhms
	s.MaxCurrentA = params.INA219.MaxCurrentA
	s.MaxShuntVoltageV = params.INA219.MaxShuntVoltageV
	s.MaxBusVoltageV = params.INA219.MaxBusVoltageV
	s.CurrentResolutionMilliAmps = params.INA219.CurrentResolutionMilliAmps

	// Original config value from Adafruit C++ code:
	// s.Configuration = configBusVoltageRange16V | configGain8_320mV | configBusADCRes12Bit |
	// 	configShuntADCRes12Bit128Sample69ms | configModeShuntAndBusVoltageContinuous

	// set up the configuration parameter
	s.Configuration =
		configBusVoltageRange32V |
			configGain1_40mV |
			configBusADCRes12Bit |
			configShuntADCRes12Bit128Sample69ms |
			configModeShuntAndBusVoltageContinuous

	calibration, err := s.calculateCalibrationValue()
	if err != nil {
		return
	}
	s.CalibrationValue = calibration
	s.currentAdj = s.CurrentResolutionMilliAmps  // Note: since current is calculated from voltage shunt, effective current resolution is 10uV/Rshunt
	s.busVoltageAdj = busVoltageLSBmV / 1000     // in V (INA219 has a fixed bus voltage LSB of 4 mV)
	s.shuntVoltageAdj = shuntVoltageLSBuV / 1000 // in mV (INA219 has a fixed shunt voltage LSB of 10uV = 0.01 mV)
	s.printConfig()

	return nil
}

func (s *INA219) applyConfiguration() (err error) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, s.Configuration)
	// push the configuration to the sensor
	err = s.bus.WriteByteBlock(s.ID(), regConfig, data)
	if err != nil {
		return
	}

	binary.BigEndian.PutUint16(data, s.CalibrationValue)
	err = s.bus.WriteByteBlock(s.ID(), regCalibration, data)
	if err != nil {
		return
	}
	return
}

// Reset performs a soft reset on the sensor.
func (s *INA219) Reset() (err error) {
	s.writeCommandWithReg(splitCmd(configResetCmd))
	return s.applyConfiguration()
}

func twosComplement(rawData []byte) (result int16) {
	if rawData[0]>>7 == 1 {
		result = int16(rawData[0])*int16(256) + int16(rawData[1])
		logger.Debugf("interim res: %d", result)
		// if result&(1<<15) != 0 {
		// 	result = result - 1<<15
		// }
	} else {
		result = int16(rawData[0])<<8 | int16(rawData[1])
	}
	return
}

// BusVoltage returns the voltage drop across the shunt resistor.
func (s *INA219) BusVoltage() (voltageVolts float64, err error) {
	err = s.applyConfiguration()
	if err != nil {
		return
	}

	vBusData, err := s.bus.ReadByteBlock(s.ID(), regBusVoltage, 2)
	if err != nil {
		return 0, err
	}
	logger.Debugf("Raw register vBusData: 0x%s", hex.EncodeToString(vBusData))

	// Note: INA219 spec says that we need to shift to the right by 3 bits, then multiply by 4 to get the value in mV
	vbus := twosComplement(vBusData) >> 3
	logger.Debugf("Shifted by 3 bits: 0x%x", vbus)
	voltageVolts = float64(vbus) * s.busVoltageAdj

	return voltageVolts, nil
}

// ShuntVoltage returns the voltage drop across the shunt resistor.
func (s *INA219) ShuntVoltage() (voltageMilliVolts float64, err error) {
	err = s.applyConfiguration()
	if err != nil {
		return
	}

	vShuntData, err := s.bus.ReadByteBlock(s.ID(), regShuntVoltage, 2)
	if err != nil {
		return 0, err
	}

	logger.Debugf("Raw register vBusData: 0x%s", hex.EncodeToString(vShuntData))
	voltageMilliVolts = float64(twosComplement(vShuntData)) * s.shuntVoltageAdj

	// if the value is very low (within 1 resolution step), the real value is probably 0
	if math.Abs(voltageMilliVolts) <= shuntVoltageLSBuV/1000 {
		voltageMilliVolts = 0
	}

	return voltageMilliVolts, nil
}

// Current returns the current through the shunt resistor.
func (s *INA219) Current() (currentMilliAmps float64, err error) {
	//err = s.applyConfiguration()
	if err != nil {
		return
	}

	currentData, err := s.bus.ReadByteBlock(s.ID(), regCurrent, 2)
	if err != nil {
		return 0, err
	}

	logger.Debugf("Raw register vBusData: 0x%s", hex.EncodeToString(currentData))
	currentMilliAmps = float64(twosComplement(currentData)) * s.currentAdj

	// if the value is very low (within 1 resolution step), the real value is probably 0
	if math.Abs(currentMilliAmps) <= (shuntVoltageLSBuV / s.ShuntResistorMilliOhms) {
		currentMilliAmps = 0
	}

	return
}

// NewMeasurement returns a new set of data from the sensor representing a Measurement unit.
func (s *INA219) NewMeasurement() (m devices.Measurement, err error) {
	sm := models.ShuntMeasurement{}
	current, err := s.Current()
	if err != nil {
		return
	}
	sm.CurrentAmps = float32(current)
	busVoltage, err := s.BusVoltage()
	if err != nil {
		return
	}
	shuntVoltage, err := s.ShuntVoltage()
	if err != nil {
		return
	}
	sm.VoltageVolts = float32(busVoltage + shuntVoltage/1000)
	sm.ShuntID = int(s.ID())
	sm.IntervalSecs = 10 // Note: this will be removed later
	sm.Timestamp = time.Now()
	m = &sm
	return m, nil
}
