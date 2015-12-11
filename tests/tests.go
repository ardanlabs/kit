// Package tests provides the generic support all tests require.
package tests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/db/mongo"
	"github.com/ardanlabs/kit/log"
)

// Success is a unicode codepoint for a check mark.
var Success = "\u2713"

// Failed is a unicode codepoint for a check X mark.
var Failed = "\u2717"

// logdash is the central buffer where all logs are stored.
var logdash bytes.Buffer

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

// Init initializes the log package.
func Init(cfgKey string) {
	cfg.Init(cfgKey)
	log.Init(&logdash, func() int { return log.DEV })
}

// InitMongo initializes the mongodb connections for testing.
func InitMongo() {
	if err := mongo.InitMGO(); err != nil {
		log.Error("Test", "Init", err, "Completed")
		logdash.WriteTo(os.Stdout)
		os.Exit(1)
	}
}

// NewRequest used to setup a request for mocking API calls with httptreemux.
func NewRequest(method, path string, body io.Reader) *http.Request {
	r, _ := http.NewRequest(method, path, body)
	u, _ := url.Parse(path)
	r.URL = u
	r.RequestURI = path

	/*
		db.auth_users.insert(
			{
			    "_id" : ObjectId("5660bc6e16908cae692e0593"),
			    "public_id" : "d648d9d1-f3a7-4586-b64e-f8d61ca986fe",
			    "private_id" : "5d829805-d801-408e-b418-2e9055da244b",
			    "status" : NumberInt(1),
			    "full_name" : "TEST USER DON'T DELETE",
			    "email" : "bill@ardanstudios.com",
			    "password" : "$2a$10$CRoh/8Uex49hviQYDlDvruoQUO10QxVOU7O0UMliqGlXSySK4SZEq",
			    "is_deleted" : false,
			    "date_modified" : ISODate("2015-12-03T22:04:30.117+0000"),
			    "date_created" : ISODate("2015-12-03T22:04:30.117+0000")
			}
		)

		db.sessions.insert(
			{
			    "_id" : ObjectId("5660bc6e16908cae692e0594"),
			    "session_id" : "6d72e6dd-93d0-4413-9b4c-8546d4d3514e",
			    "public_id" : "d648d9d1-f3a7-4586-b64e-f8d61ca986fe",
			    "date_expires" : ISODate("2016-12-02T22:04:30.282+0000"),
			    "date_created" : ISODate("2015-12-03T22:04:30.282+0000")
			}
		)
	*/

	// Add header for authentication.
	r.Header.Set("Authorization", "Basic NmQ3MmU2ZGQtOTNkMC00NDEzLTliNGMtODU0NmQ0ZDM1MTRlOlBDeVgvTFRHWjhOdGZWOGVReXZObkpydm4xc2loQk9uQW5TNFpGZGNFdnc9")

	return r
}
