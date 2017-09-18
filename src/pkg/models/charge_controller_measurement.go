package models

import (
	"fmt"
	"pkg/devices"
	"pkg/logger"
	"time"
)

// ChargeControllerMeasurement contains data obtained from a charge controller.
type ChargeControllerMeasurement struct {
	// Important: field names as well as json attributes need to be identical to those of gotracer.TracerStatus
	// to allow marshalling from gotracer.TracerStatus to ChargeControllerMeasurement
	ArrayVoltage float32 `json:"pvv" db:"ArrayVoltage_V"` // Solar panel voltage, (V)
	ArrayCurrent float32 `json:"pvc" db:"ArrayCurrent_A"` // Solar panel current, (A)
	//	ArrayPower             float32   `json:"pvp"`                                  // Solar panel power, (W)
	BatteryVoltage    float32 `json:"bv" db:"BatteryVoltage_V"`             // Battery voltage, (V)
	BatteryCurrent    float32 `json:"bc" db:"BatteryCurrent_A"`             // Battery current, (A)
	BatterySOC        int32   `json:"bsoc" db:"BatterySOC"`                 // Battery state of charge, (%)
	BatteryTemp       float32 `json:"btemp" db:"BatteryTemperature_C"`      // Battery temperatur, (C)
	BatteryMaxVoltage float32 `json:"bmaxv" db:"BatteryMinVoltage_V"`       // Battery maximum voltage, (V)
	BatteryMinVoltage float32 `json:"bminv" db:"BatteryMaxVoltage_V"`       // Battery lowest voltage, (V)
	DeviceTemp        float32 `json:"devtemp" db:"ControllerTemperature_C"` // Tracer temperature, (C)
	LoadVoltage       float32 `json:"lv" db:"LoadVoltage_V"`                // Load voltage, (V)
	LoadCurrent       float32 `json:"lc" db:"LoadCurrent_A"`                // Load current, (A)
	//	LoadPower              float32   `json:"lp" db:""`                             // Load power, (W)
	Load                   bool      `json:"load" db:"LoadOn"`                    // Shows whether load is on or off
	EnergyConsumedDaily    float32   `json:"ecd" db:"EnergyConsumedDaily_kWh"`    // Tracer calculated daily consumption, (kWh)
	EnergyConsumedMonthly  float32   `json:"ecm" db:"EnergyConsumedMonthly_kWh"`  // Tracer calculated monthly consumption, (kWh)
	EnergyConsumedAnnual   float32   `json:"eca" db:"EnergyConsumedAnnual_kWh"`   // Tracer calculated annual consumption, (kWh)
	EnergyConsumedTotal    float32   `json:"ect" db:"EnergyConsumedTotal_kWh"`    // Tracer calculated total consumption, (kWh)
	EnergyGeneratedDaily   float32   `json:"egd" db:"EnergyGeneratedDaily_kWh"`   // Tracer calculated daily power generation, (kWh)
	EnergyGeneratedMonthly float32   `json:"egm" db:"EnergyGeneratedMonthly_kWh"` // Tracer calculated monthly power generation, (kWh)
	EnergyGeneratedAnnual  float32   `json:"ega" db:"EnergyGeneratedAnnual_kWh"`  // Tracer calculated annual power generation, (kWh)
	EnergyGeneratedTotal   float32   `json:"egt" db:"EnergyGeneratedTotal_kWh"`   // Tracer calculated total power generation, (kWh)
	ID                     int       `json:"id" db:"Id"`
	Timestamp              time.Time `json:"t" db:"Timestamp"`
	ChargeControllerID     int       `json:"ccId" db:"ChargeControllerId"`
}

func (m *ChargeControllerMeasurement) String() string {
	s := fmt.Sprintf("Timestamp: %s\n", m.Timestamp)
	s += fmt.Sprintf("ID: %d\n", m.ID)
	s += fmt.Sprintf("ChargeController ID: %d\n", m.ChargeControllerID)
	s += fmt.Sprintf("ArrayVoltage: %.2f V\n", m.ArrayVoltage)
	s += fmt.Sprintf("ArrayCurrent: %.2f A\n", m.ArrayCurrent)
	s += fmt.Sprintf("BatteryVoltage: %.2f V\n", m.BatteryVoltage)
	s += fmt.Sprintf("BatteryCurrent: %.2f A\n", m.BatteryCurrent)
	s += fmt.Sprintf("BatterySOC: %d %%%%\n", m.BatterySOC)
	s += fmt.Sprintf("BatteryTemp: %.2f C\n", m.BatteryTemp)
	s += fmt.Sprintf("BatteryMaxVoltage: %.2f V\n", m.BatteryMaxVoltage)
	s += fmt.Sprintf("BatteryMinVoltage: %.2f V\n", m.BatteryMinVoltage)
	s += fmt.Sprintf("DeviceTemp: %.2f C\n", m.DeviceTemp)
	s += fmt.Sprintf("LoadVoltage: %.2f V\n", m.LoadVoltage)
	s += fmt.Sprintf("LoadCurrent: %.2f C\n", m.LoadCurrent)
	s += fmt.Sprintf("Load on: %t\n", m.Load)
	s += fmt.Sprintf("EnergyConsumedDaily: %.2f KWh\n", m.EnergyConsumedDaily)
	s += fmt.Sprintf("EnergyConsumedMonthly: %.2f KWh\n", m.EnergyConsumedMonthly)
	s += fmt.Sprintf("EnergyConsumedAnnual: %.2f KWh\n", m.EnergyConsumedAnnual)
	s += fmt.Sprintf("EnergyConsumedTotal: %.2f KWh\n", m.EnergyConsumedTotal)

	s += fmt.Sprintf("EnergyGeneratedDaily: %.2f KWh\n", m.EnergyGeneratedDaily)
	s += fmt.Sprintf("EnergyGeneratedMonthly: %.2f KWh\n", m.EnergyGeneratedMonthly)
	s += fmt.Sprintf("EnergyGeneratedAnnual: %.2f KWh\n", m.EnergyGeneratedAnnual)
	s += fmt.Sprintf("EnergyGeneratedTotal: %.2f KWh", m.EnergyGeneratedTotal)

	return s
}

// InsertIntoDb inserts the measurement into the database.
func (m *ChargeControllerMeasurement) InsertIntoDb(dbmap devices.DbInserter) error {
	logger.Debugf("ChargeControllerMeasurement.InsertIntoDb")

	var m2 devices.Measurement
	m2 = m

	return dbmap.Insert(m2)

	//return dbmap.Insert(m)
}
