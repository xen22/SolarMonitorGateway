package helpers

import (
	"cmd/solarcmd/config"
	"pkg/logger"
	"time"

	"github.com/spagettikod/gotracer"
)

// QueryChargeController retrieves data from the charge controller
func QueryChargeController() (status gotracer.TracerStatus, err error) {
	if tracerAvailable {
		logger.Info("Getting charge controller data from Tracer4215BN.")
		status, err = gotracer.Status(config.TracerSerialDevice)
		if err != nil {
			logger.Fatalf("Unable to get status from Tracer serial device %q. Err, %s.", config.TracerSerialDevice, err)
			return
		}
	} else {
		status = getSampleTracerStatus()
	}
	status.Timestamp = time.Now()

	logger.Infof("Charge controller data: \n%s", status.String())
	return
}

// ==================================================================================================================
// Note: this is for testing only (to be removed)
// variables used for testing when actual devices are not available
const tracerAvailable = false

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

// ==================================================================================================================
