package solar

import (
	"cmd/solarcmd/config"
	"pkg/devices"
	"pkg/logger"
	"time"

	"pkg/adapters"

	"github.com/spagettikod/gotracer"
)

// Tracer4215BN implements a EPSolar Tracer 4215BN charge controller.
type Tracer4215BN struct {
	devices.DeviceInfo
	SerialDevice string
}

// NewTracer4215BN creates a new Tracer4215BN object and initialises it.
func NewTracer4215BN(config devices.DeviceConfig) (t *Tracer4215BN, err error) {
	t = &Tracer4215BN{}
	err = t.Init(config)
	return
}

// Init is used to initialise an Tracer4215BN object.
func (t *Tracer4215BN) Init(config devices.DeviceConfig) (err error) {
	err = t.DeviceInfo.Init(config)
	if err != nil {
		return
	}

	t.SerialDevice = config.Params.Tracer.SerialDevice
	return nil
}

// NewMeasurement retrieves data from the charge controller.
func (t *Tracer4215BN) NewMeasurement() (dm devices.Measurement, err error) {

	status := gotracer.TracerStatus{}
	logger.Info("Getting charge controller data from Tracer4215BN.")
	status, err = gotracer.Status(config.TracerSerialDevice)
	if err != nil {
		logger.Fatalf("Unable to get status from Tracer serial device %q. Err, %s.", config.TracerSerialDevice, err)
		return
	}
	status.Timestamp = time.Now()
	//	logger.Infof("Charge controller data: \n%s", status.String())
	m, err := adapters.ChargeControllerMeasurementFromTracerStatus(status)
	if err != nil {
		return
	}
	// save database ID of parent ChargeController
	m.ChargeControllerID = int(t.ID())
	dm = &m

	return
}
