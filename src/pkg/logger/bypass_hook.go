package logger

import "github.com/Sirupsen/logrus"

// BypassHook is used to bypass the logging system and just print the data to the screen without formatting.
type BypassHook struct {
}

// Levels returns the levels to which this hook applies.
func (hook *BypassHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel, logrus.FatalLevel}
}

// Fire is used to perform the work of the hook.
func (hook *BypassHook) Fire(entry *logrus.Entry) error {
	entry.Data = logrus.Fields{}
	//	fmt.Println(entry.Message)
	return nil
}
