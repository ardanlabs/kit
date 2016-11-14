// Package tests provides the generic support all tests require.
package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// Context provides a base context for tests.
var Context = "Test"

// TestSession is the name used to register the MongoDB session.
var TestSession = "test"

// Success and failure markers.
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

// logdash is the central buffer where all logs are stored.
var logdash bytes.Buffer

//==============================================================================

// ResetLog resets the contents of logdash.
func ResetLog() {
	logdash.Reset()
}

// DisplayLog writes the logdash data to standand out, if testing in verbose mode
// was turned on.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}

	logdash.WriteTo(os.Stdout)
}

// IndentJSON takes a JSON payload as a string and re-indents it to make
// comparing expected strings to tests strings during testing.
func IndentJSON(j string) string {
	var indented interface{}
	if err := json.Unmarshal([]byte(j), &indented); err != nil {
		return ""
	}

	data, err := json.MarshalIndent(indented, "", "  ")
	if err != nil {
		return ""
	}

	return string(data)
}
