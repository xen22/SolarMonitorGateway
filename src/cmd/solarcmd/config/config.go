package config

import "pkg/devices"

// Config holds the configuration data
type Config struct {
	Database devices.DatabaseConfig `json:"Database"`
	Sensors  []devices.DeviceConfig `json:"Sensors"`
}
