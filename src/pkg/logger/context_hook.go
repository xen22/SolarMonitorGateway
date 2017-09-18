package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

// ContextHook is used to print source file, line number and function on every log line.
type ContextHook struct {
}

// Levels returns the levels to which this hook applies.
func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire is used to perform the work of the hook.
// from https://github.com/Gurpartap/logrus-stack/blob/master/logrus-stack-hook.go
// func (hook *ContextHook) Fire(entry *logrus.Entry) error {
// 	var skipFrames int
// 	if len(entry.Data) == 0 {
// 		// When WithField(s) is not used, we have 8 logrus frames to skip.
// 		skipFrames = 8
// 	} else {
// 		// When WithField(s) is used, we have 6 logrus frames to skip.
// 		skipFrames = 6
// 	}

// 	var frames stack.Stack

// 	// Get the complete stack track past skipFrames count.
// 	_frames := stack.Callers(skipFrames)

// 	// Remove logrus's own frames that seem to appear after the code is through
// 	// certain hoops. e.g. http handler in a separate package.
// 	// This is a workaround.
// 	for _, frame := range _frames {
// 		if !strings.Contains(frame.File, "github.com/Sirupsen/logrus") {
// 			frames = append(frames, frame)
// 		}
// 	}

// 	if len(frames) > 0 {
// 		// If we have a frame, we set it to "caller" field for assigned levels.
// 		for _, level := range hook.CallerLevels {
// 			if entry.Level == level {
// 				entry.Data["caller"] = frames[0]
// 				break
// 			}
// 		}

// 		// Set the available frames to "stack" field.
// 		for _, level := range hook.StackLevels {
// 			if entry.Level == level {
// 				entry.Data["stack"] = frames
// 				break
// 			}
// 		}
// 	}

// 	return nil
// }

// Fire is used to perform the work of the hook.
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(8, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/Sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			// entry.Data["file"] = path.Base(file)
			// entry.Data["func"] = path.Base(name)
			// entry.Data["line"] = line
			entry.Data["f"] = fmt.Sprintf("%s() [%s:%d]", path.Base(name), path.Base(file), line)
			break
		}
	}
	return nil
}
