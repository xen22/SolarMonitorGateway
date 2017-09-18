package sim

import (
	"pkg/adapters"
	"pkg/devices"
	"pkg/logger"
	"pkg/solar"
	"time"

	"github.com/spagettikod/gotracer"
)

// SimulatedTracer4215BN simulates a Tracer4215BN charge controller
type SimulatedTracer4215BN struct {
	solar.Tracer4215BN
}

// NewSimulatedTracer4215BN creates a new Tracer4215BN object and initialises it.
func NewSimulatedTracer4215BN(config devices.DeviceConfig) (t *SimulatedTracer4215BN, err error) {
	t = &SimulatedTracer4215BN{}
	err = t.Init(config)
	return
}

// NewMeasurement generates a new measurement to simulate retrieval of data from the controller.
func (t *SimulatedTracer4215BN) NewMeasurement() (dm devices.Measurement, err error) {

	status := gotracer.TracerStatus{}
	logger.Info("Simulating getting charge controller data from Tracer4215BN.")
	status = getSampleTracerStatus()
	status.Timestamp = time.Now()
	m, err := adapters.ChargeControllerMeasurementFromTracerStatus(status)
	if err != nil {
		return
	}
	// save database ID of parent ChargeController
	m.ChargeControllerID = int(t.ID())
	dm = &m
	return
}

func getSampleTracerStatus() gotracer.TracerStatus {
	status := gotracer.TracerStatus{
		ArrayVoltage:           68.46,
		ArrayCurrent:           0.07,
		ArrayPower:             4.89,
		BatteryVoltage:         13.59,
		BatteryCurrent:         0.37,
		BatterySOC:             100,
		BatteryTemp:            19.13,
		BatteryMaxVoltage:      16.31,
		BatteryMinVoltage:      13.59,
		DeviceTemp:             17.78,
		LoadVoltage:            0.00,
		LoadCurrent:            0.00,
		LoadPower:              0.00,
		Load:                   false,
		EnergyConsumedDaily:    0.00,
		EnergyConsumedMonthly:  0.00,
		EnergyConsumedAnnual:   0.00,
		EnergyConsumedTotal:    0.00,
		EnergyGeneratedDaily:   0.03,
		EnergyGeneratedMonthly: 0.34,
		EnergyGeneratedAnnual:  34.48,
		EnergyGeneratedTotal:   34.48,
	}
	return status
}
