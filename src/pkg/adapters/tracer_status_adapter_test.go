package adapters

import (
	"testing"
	"time"

	"github.com/spagettikod/gotracer"
)

func TestChargeControllerMeasurementFromTracerStatus(t *testing.T) {

	s := gotracer.TracerStatus{}
	s.ArrayVoltage = 2.0
	s.ArrayCurrent = 1.3
	s.BatteryVoltage = 4.33
	s.BatteryCurrent = 11
	s.BatterySOC = 40
	s.BatteryTemp = 50.0005
	s.BatteryMaxVoltage = 221
	s.BatteryMinVoltage = 100
	s.DeviceTemp = 1
	s.LoadVoltage = 30
	s.LoadCurrent = -4.55
	s.Load = true
	s.EnergyConsumedDaily = 20.2
	s.EnergyConsumedMonthly = 1011.1
	s.EnergyConsumedAnnual = 122.0
	s.EnergyConsumedTotal = 10
	s.EnergyGeneratedDaily = 30
	s.EnergyGeneratedMonthly = 40
	s.EnergyGeneratedAnnual = 51.1
	s.EnergyGeneratedTotal = 14.1
	s.Timestamp = time.Now()

	m, err := ChargeControllerMeasurementFromTracerStatus(s)

	if err != nil {
		t.Errorf("Converting gotracer.TracerStatus to ChargeControllerMeasurement failed. Error returned: %s", err)
	}

	if m.ArrayVoltage != s.ArrayVoltage ||
		m.ArrayCurrent != s.ArrayCurrent ||
		m.BatteryVoltage != s.BatteryVoltage ||
		m.BatteryCurrent != s.BatteryCurrent ||
		m.BatterySOC != s.BatterySOC ||
		m.BatteryTemp != s.BatteryTemp ||
		m.BatteryMaxVoltage != s.BatteryMaxVoltage ||
		m.BatteryMinVoltage != s.BatteryMinVoltage ||
		m.DeviceTemp != s.DeviceTemp ||
		m.LoadVoltage != s.LoadVoltage ||
		m.LoadCurrent != s.LoadCurrent ||
		m.Load != s.Load ||
		m.EnergyConsumedDaily != s.EnergyConsumedDaily ||
		m.EnergyConsumedMonthly != s.EnergyConsumedMonthly ||
		m.EnergyConsumedAnnual != s.EnergyConsumedAnnual ||
		m.EnergyConsumedTotal != s.EnergyConsumedTotal ||
		m.EnergyGeneratedDaily != s.EnergyGeneratedDaily ||
		m.EnergyGeneratedMonthly != s.EnergyGeneratedMonthly ||
		m.EnergyGeneratedAnnual != s.EnergyGeneratedAnnual ||
		m.EnergyGeneratedTotal != s.EnergyGeneratedTotal ||
		m.Timestamp != s.Timestamp {
		t.Errorf("The returned ChargeControllerMeasurement object is not the same as the TracerStatus object passed in.")
	}

	if m.ID != 0 || m.ChargeControllerID != 0 {
		t.Errorf("The ChargeControllerMeasurement IDs are not initialised correctly.")
	}
}
