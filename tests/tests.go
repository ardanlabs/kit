// Package tests provides the generic support all tests require.
package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// TraceID provides a base trace id for tests.
var TraceID = "Test"

// Success and failure markers.
var (
	Success = "\u2713"
	Failed  = "\u2717"
)

// Logdash is the central buffer where all logs are stored.
var Logdash bytes.Buffer

// ResetLog resets the contents of Logdash.
func ResetLog() {
	Logdash.Reset()
}

// DisplayLog writes the Logdash data to standand out, if testing in verbose mode
// was turned on.
func DisplayLog() {
	if !testing.Verbose() {
		return
	}

	Logdash.WriteTo(os.Stdout)
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
