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

const (
	// Ldate enables the date in the local time zone: 2009/01/23
	Ldate = 1 << iota
	// Ltime enables the time in the local time zone: 01:23:23
	Ltime
	// Lmicroseconds enables microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Lmicroseconds
	// Llongfile enables full file name and line number: /a/b/c/d.go:23
	Llongfile
	// Lshortfile enables final file name element and line number: d.go:23. overrides Llongfile
	Lshortfile
	// LUTC enables if Ldate or Ltime is set, use UTC rather than the local time zone
	LUTC
	// LstdFlags enables initial values for the standard logger
	LstdFlags = Ldate | Ltime
	// Ldefault enables intial values for the default kit logger
	Ldefault = log.Ldate | log.Ltime | log.Lshortfile
)

// Logger contains a standard logger for all logging.
type Logger struct {
	*log.Logger
	level func() int
	mu    sync.RWMutex
}

//==============================================================================

// New returns a instance of a logger.
func New(w io.Writer, levelHandler func() int, flags int) *Logger {
	return &Logger{
		Logger: log.New(w, "", flags),
		level:  levelHandler,
	}
}

//==============================================================================

// mLevel sets the default log level for use with the log methods.
const mLevel = 2

// Dev logs trace information for developers.
func (l *Logger) Dev(context interface{}, funcName string, format string, a ...interface{}) {
	l.mu.RLock()
	{
		if l.level() == DEV {
			if a != nil {
				format = fmt.Sprintf(format, a...)
			}

			l.Output(mLevel, fmt.Sprintf("DEV : %s : %s : %s", context, funcName, format))
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

			l.Output(mLevel, fmt.Sprintf("USER : %s : %s : %s", context, funcName, format))
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

			l.Output(mLevel, fmt.Sprintf("ERROR : %s : %s : %s : %s", context, funcName, err, format))
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

			l.Output(mLevel, fmt.Sprintf("FATAL : %s : %s : %s", context, funcName, format))
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

			l.Output(mLevel+offset, fmt.Sprintf("DEV : %s : %s : %s", context, funcName, format))
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

			l.Output(mLevel+offset, fmt.Sprintf("USER : %s : %s : %s", context, funcName, format))
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

			l.Output(mLevel+offset, fmt.Sprintf("ERROR : %s : %s : %s : %s", context, funcName, err, format))
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

			l.Output(mLevel+offset, fmt.Sprintf("FATAL : %s : %s : %s", context, funcName, format))
		}
	}
	l.mu.RUnlock()

	os.Exit(1)
}
