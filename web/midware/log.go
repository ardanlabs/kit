package midware

import (
	"os"

	"github.com/ardanlabs/kit/cfg"
	kitlog "github.com/ardanlabs/kit/log"
)

//==============================================================================

// logger provides the internal midware package log interface.
type logger interface {
	Dev(interface{}, string, string, ...interface{})
	User(interface{}, string, string, ...interface{})
	Fatal(interface{}, string, string, ...interface{})
	Error(interface{}, string, error, string, ...interface{})
}

//==============================================================================

// cfgLogLevel defines the key to retrieve the log level for the logger.
var cfgLogLevel = "LOGGING_LEVEL"

// log provides the default log instance for the midware package.
var log logger

func init() {
	level, err := cfg.Int(cfgLogLevel)
	if err != nil {
		level = kitlog.DEV
	}

	log = kitlog.New(os.Stdout, func() int { return level })
}

//==============================================================================
