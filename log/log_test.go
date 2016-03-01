package log_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ardanlabs/kit/log"
)

// Success and failure markers.
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

// logdest implements io.Writer and is the log package destination.
var logdest bytes.Buffer

// resetLog can be called at the beginning of a test or example.
func resetLog() { logdest.Reset() }

// displayLog can be called at the end of a test or example.
// It only prints the log contents if the -test.v flag is set.
func displayLog() {
	if !testing.Verbose() {
		return
	}
	logdest.WriteTo(os.Stdout)
}

//==============================================================================

// TestLogLevelUSER tests the basic functioning of the logger in USER mode.
func TestLogLevelUSER(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to USER.")
		{
			log.Init(&logdest, func() int { return log.USER })
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:53: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:54: ERROR : context : FuncName : An error : Message 3 no format\n", dt)

			log.Dev("context", "FuncName", "Message 1 no format")
			log.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")
			log.Error("context", "FuncName", errors.New("An error"), "Message 3 no format")

			if logdest.String() == log1+log2 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}

// TestLogLevelDEV tests the basic functioning of the logger in DEV mode.
func TestLogLevelDEV(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to DEV.")
		{
			log.Init(&logdest, func() int { return log.DEV })
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:83: DEV : context : FuncName : Message 1 no format\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:84: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log3 := fmt.Sprintf("%s log_test.go:85: ERROR : context : FuncName : An error : Message 3 with format: C, D\n", dt)

			log.Dev("context", "FuncName", "Message 1 no format")
			log.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")
			log.Error("context", "FuncName", errors.New("An error"), "Message 3 with format: %s, %s", "C", "D")

			if logdest.String() == log1+log2+log3 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2+log3)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}

// TestLogInstanceInDev tests the basic functioning of the logger instance
func TestLogInstanceInDev(t *testing.T) {
	t.Log("Given the need to log DEV and DEV messages.")
	{
		t.Log("\tWhen we set the logging level to DEV.")
		{
			lg := log.New(&logdest, func() int { return log.DEV })

			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:115: DEV : context : FuncName : Message 1 no format\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:116: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log3 := fmt.Sprintf("%s log_test.go:117: ERROR : context : FuncName : An error : Message 3 with format: C, D\n", dt)

			lg.Dev("context", "FuncName", "Message 1 no format")
			lg.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")
			lg.Error("context", "FuncName", errors.New("An error"), "Message 3 with format: %s, %s", "C", "D")

			if logdest.String() == log1+log2+log3 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2+log3)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}

// TestLogInstanceInUser tests the basic functioning of the logger instance
func TestLogInstanceInUser(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to USER.")
		{
			// log.Init(&logdest, func() int { return log.USER })
			lg := log.New(&logdest, func() int { return log.USER })

			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:148: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:149: ERROR : context : FuncName : An error : Message 3 no format\n", dt)

			lg.Dev("context", "FuncName", "Message 1 no format")
			lg.User("context", "FuncName", "Message 2 with format: %s, %s", "A", "B")
			lg.Error("context", "FuncName", errors.New("An error"), "Message 3 no format")

			if logdest.String() == log1+log2 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}

// TestLogLevelDEVOffset tests the basic functioning of the logger in DEV mode.
func TestLogLevelDEVOffset(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to DEV.")
		{
			log.Init(&logdest, func() int { return log.DEV })
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:178: DEV : context : FuncName : Message 1 no format\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:179: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log3 := fmt.Sprintf("%s log_test.go:180: ERROR : context : FuncName : An error : Message 3 with format: C, D\n", dt)

			log.DevOffset("context", 0, "FuncName", "Message 1 no format")
			log.UserOffset("context", 0, "FuncName", "Message 2 with format: %s, %s", "A", "B")
			log.ErrorOffset("context", 0, "FuncName", errors.New("An error"), "Message 3 with format: %s, %s", "C", "D")

			if logdest.String() == log1+log2+log3 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2+log3)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}

// TestLogLevelUserOffset tests the basic functioning of the logger in DEV mode.
func TestLogLevelUserOffset(t *testing.T) {
	t.Log("Given the need to log DEV and USER messages.")
	{
		t.Log("\tWhen we set the logging level to DEV.")
		{
			log.Init(&logdest, func() int { return log.USER })
			resetLog()
			defer displayLog()

			dt := time.Now().Format("2006/01/02 15:04:05")

			log1 := fmt.Sprintf("%s log_test.go:209: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:210: ERROR : context : FuncName : An error : Message 3 with format: C, D\n", dt)

			log.DevOffset("context", 0, "FuncName", "Message 1 no format")
			log.UserOffset("context", 0, "FuncName", "Message 2 with format: %s, %s", "A", "B")
			log.ErrorOffset("context", 0, "FuncName", errors.New("An error"), "Message 3 with format: %s, %s", "C", "D")

			if logdest.String() == log1+log2 {
				t.Logf("\t\t%v : Should log the expected trace line.", Success)
			} else {
				t.Log("***>", logdest.String())
				t.Log("***>", log1+log2)
				t.Errorf("\t\t%v : Should log the expected trace line.", Failed)
			}
		}
	}
}
