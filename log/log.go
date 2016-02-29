package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

//==============================================================================

// Log provides a log interface for the log package.
// type Log interface {
// 	Dev(interface{}, string, string, ...interface{})
// 	User(interface{}, string, string, ...interface{})
// 	Fatal(interface{}, string, string, ...interface{})
// 	Error(interface{}, string, string, ...interface{})
// }

//==============================================================================

// Level constants that define the supported usable LogLevel.
const (
	NONE int = iota
	DEV
	USER
)

//==============================================================================

// Logger contains a standard logger for all logging.
type Logger struct {
	*log.Logger
	level func() int
	mu    sync.RWMutex
}

// New returns a instance of a logger.
func New(w io.Writer, level func() int) *Logger {
	lm := Logger{
		Logger: log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile),
		level:  level,
	}

	return &lm
}

// Dev logs trace information for developers.
func (l *Logger) Dev(context interface{}, funcName string, format string, a ...interface{}) {
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
func (l *Logger) User(context interface{}, funcName string, format string, a ...interface{}) {
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
func (l *Logger) Error(context interface{}, funcName string, err error, format string, a ...interface{}) {
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
func (l *Logger) Fatal(context interface{}, funcName string, format string, a ...interface{}) {
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

//==============================================================================

// l defines the default log variable for the global log functions.
var l Logger

// Init initializes the default logger to allow usage of the global log
// functions.
func Init(w io.Writer, level func() int) {
	l.mu.Lock()
	{
		l.Logger = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
		l.level = level
	}
	l.mu.Unlock()
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

//==============================================================================
