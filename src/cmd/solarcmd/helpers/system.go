package helpers

import (
	"cmd/solarcmd/config"
	"pkg/db"
	"pkg/devices"
	"pkg/factories"
	"pkg/logger"
	"strings"
)

// System is the main type with functions used to retrieve measurements from the system's devices.
type System struct {
	devices []devices.Device
	factory *factories.DeviceFactory
	config  config.Config
}

// NewSystem creates a new System object.
func NewSystem(configFile string) (s *System, err error) {
	s = &System{}
	err = s.init(configFile)
	return
}

func (s *System) init(configFile string) (err error) {
	s.config, err = config.ParseConfig(configFile)
	if err != nil {
		return
	}

	logger.Debugf("config data: %+v", s.config)

	s.factory = &factories.DeviceFactory{}
	for _, devConfig := range s.config.Sensors {
		d, err := s.factory.Create(devConfig)
		if err != nil {
			return err
		}
		s.devices = append(s.devices, d)
	}
	return err
}

// NewMeasurements retrieves a list of measurements from configured devices.
func (s *System) NewMeasurements() (mlist []devices.Measurement, err error) {
	mlist = make([]devices.Measurement, len(s.devices))
	for _, d := range s.devices {
		m, err := d.NewMeasurement()
		if err != nil {
			return nil, err
		}
		mlist = append(mlist, m)
	}
	return
}

// PrintMeasurements prints measurements from configured devices.
func (s *System) PrintMeasurements() (err error) {
	logger.Debugf("PrintMeasurements")
	for _, d := range s.devices {
		logger.Debugf("Device: %s", d.Name())
		m, err := d.NewMeasurement()
		if err != nil {
			return err
		}
		logger.Infof("-----> %s:", d.Name())
		for _, line := range strings.Split(m.String(), "\n") {
			logger.Infof(line)
		}
	}
	return nil
}

// SaveMeasurementsToDb saves measurements to the database.
func (s *System) SaveMeasurementsToDb() (err error) {
	db := db.NewSolarDb(s.config.Database)
	err = db.Init()
	if err != nil {
		return err
	}
	defer func() {
		logger.Info("Closing db connection.")
		db.Close()
	}()

	for _, d := range s.devices {
		m, err := d.NewMeasurement()
		if err != nil {
			return err
		}
		logger.Debugf("SaveMeasurementsToDb: %s: %s", d.Name, m)
		err = db.InsertMeasurement(m)
		if err != nil {
			return err
		}
	}
	return nil
}
