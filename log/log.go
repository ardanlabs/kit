package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Level constants that define the supported usable LogLevel.
const (
	NONE int = iota
	DEV
	USER
)

// l contains a standard logger for all logging.
var l struct {
	*log.Logger
	level func() int
	mu    sync.RWMutex
}

//==============================================================================

// Init must be called to initialize the logging system. This function should
// only be called once.
func Init(w io.Writer, level func() int) error {
	l.mu.Lock()
	{
		l.Logger = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
		l.level = level
	}
	l.mu.Unlock()

	return nil
}

// Dev logs trace information for developers.
func Dev(context interface{}, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() == DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(2, fmt.Sprintf("DEV : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()
}

// User logs trace information for users.
func User(context interface{}, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(2, fmt.Sprintf("USER : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()
}

// Error logs trace information that are errors.
func Error(context interface{}, funcName string, err error, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(2, fmt.Sprintf("ERROR : %s : %s : %s : %s", context, funcName, err, format))
		}
	}
	l.mu.RUnlock()
}

// Fatal logs trace information for users and terminates the app.
func Fatal(context interface{}, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(2, fmt.Sprintf("FATAL : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()

	os.Exit(1)
}
