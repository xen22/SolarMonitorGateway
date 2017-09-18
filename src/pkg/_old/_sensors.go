package helpers

import (
	"cmd/solarcmd/config"
	"pkg/devices"
	"pkg/i2c"
	"pkg/logger"
)

// CreateSensors creates all sensors configured in the system
func CreateSensors() (tSensors map[string]devices.TemperatureSensor,
	hSensors map[string]devices.HumiditySensor, shunts map[string]devices.Shunt, err error) {

	tSensors = map[string]devices.TemperatureSensor{}
	hSensors = map[string]devices.HumiditySensor{}
	shunts = map[string]devices.Shunt{}

	sensorsInfo := config.GetSensorConfig()
	for name, sensorData := range sensorsInfo {
		switch sensorData.Type {
		case devices.TemperatureType:
			s, err := i2c.NewTemperatureSensor(sensorData)
			logger.FatalErrf(err, "Failed to create temperature sensor.")
			tSensors[name] = s

		case devices.HumidityType:
			s, err := i2c.NewHumiditySensor(sensorData)
			logger.FatalErrf(err, "Failed to create humidity sensor.")
			hSensors[name] = s

		case devices.ShuntType:
			s, err := i2c.NewShuntSensor(sensorData)
			logger.FatalErrf(err, "Failed to create shunt sensor.")
			shunts[name] = s

		default:
			logger.Infof("Ignoring sensor %q of type %d", name, sensorData.Type)
		}
	}

	return
}

// QuerySensors retrieves data from all sensors in the system
func QuerySensors() {
	logger.Info("Querying known sensors: ")

	tSensors, hSensors, shunts, err := CreateSensors()
	logger.FatalErrf(err, "Failed to create sensors")

	for name, s := range tSensors {
		logger.Debugf("Resetting sensor %s", name)
		err := s.Reset()
		logger.FatalErrf(err, "Failed to reset sensor.")

		temp, err := s.Temperature()
		logger.FatalErrf(err, "Failed to read temperature.")
		logger.Infof("[%s] Temperature: \t%.2f deg C", name, temp)
	}

	for name, s := range hSensors {
		logger.Debugf("Resetting sensor %s", name)
		err := s.Reset()
		logger.FatalErrf(err, "Failed to reset sensor.")

		humidity, err := s.Humidity()
		logger.FatalErrf(err, "Failed to read humidity.")
		logger.Infof("[%s] Humidity: \t%.2f %%", name, humidity)
	}

	for name, s := range shunts {
		logger.Debugf("Resetting sensor %s", name)
		err := s.Reset()
		logger.FatalErrf(err, "Failed to reset sensor.")

		busVoltage, err := s.BusVoltage()
		logger.FatalErrf(err, "Failed to read bus voltage.")
		logger.Infof("[%s] Bus voltage: \t%.4f V", name, busVoltage)

		shuntVoltage, err := s.ShuntVoltage()
		logger.FatalErrf(err, "Failed to read shunt voltage.")
		logger.Infof("[%s] Shunt voltage: \t%.4f mV", name, shuntVoltage)

		logger.Infof("[%s] Total voltage: \t%.4f V", name, busVoltage+shuntVoltage/1000)

		current, err := s.Current()
		logger.FatalErrf(err, "Failed to read current.")
		logger.Infof("[%s] Current: \t\t%.4f mA", name, current)
	}

}
