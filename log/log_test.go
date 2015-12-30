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

			log1 := fmt.Sprintf("%s log_test.go:51: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:52: ERROR : context : FuncName : An error : Message 3 no format\n", dt)

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

			log1 := fmt.Sprintf("%s log_test.go:81: DEV : context : FuncName : Message 1 no format\n", dt)
			log2 := fmt.Sprintf("%s log_test.go:82: USER : context : FuncName : Message 2 with format: A, B\n", dt)
			log3 := fmt.Sprintf("%s log_test.go:83: ERROR : context : FuncName : An error : Message 3 with format: C, D\n", dt)

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
