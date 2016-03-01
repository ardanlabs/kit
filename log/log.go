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
// 	DevOffset(interface{}, int,string, string, ...interface{})
// 	UserOffset(interface{}, int,string, string, ...interface{})
// 	FatalOffset(interface{}, int,string, string, ...interface{})
// 	ErrorOffset(interface{}, int,string, string, ...interface{})
// }

//==============================================================================

// Level constants that define the supported usable LogLevel.
const (
	NONE int = iota
	DEV
	USER
)

//==============================================================================

// defaultLogOffset sets the default log level for use with the log offset
// functions.
const defaultLogOffset = 2

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

			l.Output(defaultLogOffset, fmt.Sprintf("DEV : %s : %s : %s", context, funcName, format))
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

			l.Output(defaultLogOffset, fmt.Sprintf("USER : %s : %s : %s", context, funcName, format))
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

			l.Output(defaultLogOffset, fmt.Sprintf("ERROR : %s : %s : %s : %s", context, funcName, err, format))
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

			l.Output(defaultLogOffset, fmt.Sprintf("FATAL : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()

	os.Exit(1)
}

//==============================================================================

// DevOffset logs trace information for developers with a offset option to
// expand the caller level.
func (l *Logger) DevOffset(context interface{}, offset int, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() == DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(defaultLogOffset+offset, fmt.Sprintf("DEV : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()
}

// UserOffset logs trace information for users with a offset option to expand the
// caller level.
func (l *Logger) UserOffset(context interface{}, offset int, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(defaultLogOffset+offset, fmt.Sprintf("USER : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()
}

// ErrorOffset logs trace information that are errors with a offset option to
// expand the caller level.
func (l *Logger) ErrorOffset(context interface{}, offset int, funcName string, err error, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(defaultLogOffset+offset, fmt.Sprintf("ERROR : %s : %s : %s : %s", context, funcName, err, format))
		}
	}
	l.mu.RUnlock()
}

// FatalOffset logs trace information for users and terminates the app with a
// offset expand the caller level.
func (l *Logger) FatalOffset(context interface{}, offset int, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() >= DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(defaultLogOffset+offset, fmt.Sprintf("FATAL : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()

	os.Exit(1)
}
