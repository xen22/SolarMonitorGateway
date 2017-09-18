package logger

import (
	"os"
	"strings"

	"pkg/_vendor/logrus-stack"

	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type myLogger struct {
	impl    *logrus.Logger
	isDebug bool
	isQuiet bool
}

var log *myLogger

//var log *logrus.Logger

var isDebug = false
var isQuiet = false

func init() {
	log = newLogger()

}

func newLogger() *myLogger {
	ret := &myLogger{}
	ret.init()
	return ret
}

func (log *myLogger) init() {
	log.impl = logrus.New()

	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, "-debug") || strings.Contains(arg, "-d") {
			log.isDebug = true
		}
		if strings.Contains(arg, "-quiet") || strings.Contains(arg, "-q") {
			log.isQuiet = true
		}
	}

	var customFormatter logrus.Formatter
	if log.isQuiet {
		formatter := new(logrus.TextFormatter)
		formatter.DisableTimestamp = true
		//formatter.DisableColors = true
		//formatter.DisableSorting = true
		customFormatter = formatter
	} else {

		formatter := new(prefixed.TextFormatter)
		formatter.TimestampFormat = "02 Jan | 15:04:05.000"
		formatter.ShortTimestamp = false
		formatter.SpacePadding = 80
		customFormatter = formatter
	}

	log.impl.Formatter = customFormatter

	if log.isDebug {
		log.impl.Level = logrus.DebugLevel
	} else {
		log.impl.Level = logrus.InfoLevel
	}

	if log.isQuiet {
		log.impl.Hooks.Add(&BypassHook{})
	} else {
		log.impl.Hooks.Add(logrus_stack.StandardHook())
	}
}

// Info logs simple informational messages.
func Info(args ...interface{}) {
	log.impl.Info(args...)
}

// Infof logs formatted informational messages.
func Infof(format string, args ...interface{}) {
	log.impl.Infof(format, args...)
}

// Debug logs simple debug messages.
func Debug(args ...interface{}) {
	log.impl.Debug(args...)
}

// Debugf logs formatted debug messages.
func Debugf(format string, args ...interface{}) {
	log.impl.Debugf(format, args...)
}

// Fatal logs simple fatal messages.
func Fatal(args ...interface{}) {
	log.impl.Fatal(args...)
}

// Fatalf logs formatted fatal messages.
func Fatalf(format string, args ...interface{}) {
	log.impl.Fatalf(format, args...)
}

// FatalErrf logs formatted fatal messages if an error has occurred.
func FatalErrf(err error, format string, args ...interface{}) {
	if err != nil {
		log.impl.WithField("err", err).Fatalf(format, args...)
	}
}
