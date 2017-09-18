package config

import (
	"encoding/json"
	"io/ioutil"
	"pkg/logger"
)

// ParseConfig parses the json configuration file into a Config struct.
func ParseConfig(configFile string) (config Config, err error) {
	logger.Debugf("Parsing config file %s", configFile)

	buff, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(buff, &config)
	if err != nil {
		return
	}
	return
}
