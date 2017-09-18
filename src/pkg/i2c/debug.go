package i2c

import (
	"os/exec"
	"pkg/logger"
)

// RunI2CDetect executes i2cdetect for debugging purposes.
func RunI2CDetect() {
	out, err := exec.Command("i2cdetect", "-y", "1").Output()
	if err != nil {
		logger.FatalErrf(err, "Failed to execute i2cdetect.")
	}
	logger.Debug(string(out))
}
