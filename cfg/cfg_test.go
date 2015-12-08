package cfg_test

import (
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// TestExists validates the ability to load configuration values
// using the OS-level environment variables and read them back.
func TestExists(t *testing.T) {
	t.Log("Given the need to read environment variables.")
	{
		uStr := "postgres://root:root@127.0.0.1:8080/postgres?sslmode=disable"

		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")
		os.Setenv("MYAPP_FLAG", "true")
		os.Setenv("MYAPP_DSN", uStr)

		cfg.Init("MYAPP")

		t.Log("\tWhen given a namspace key to search for that exists.")
		{
			proc, err := cfg.Int("PROC_ID")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "PROC_ID")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", succeed, "PROC_ID")

				if proc != 322 {
					t.Errorf("\t\t%s Should have key %q with value %d", failed, "PROC_ID", 322)
				} else {
					t.Logf("\t\t%s Should have key %q with value %d", succeed, "PROC_ID", 322)
				}
			}

			socket, err := cfg.String("SOCKET")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "SOCKET")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", succeed, "SOCKET")

				if socket != "./tmp/sockets.po" {
					t.Errorf("\t\t%s Should have key %q with value %q", failed, "SOCKET", "./tmp/sockets.po")
				} else {
					t.Logf("\t\t%s Should have key %q with value %q", succeed, "SOCKET", "./tmp/sockets.po")
				}
			}

			port, err := cfg.Int("PORT")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "PORT")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", succeed, "PORT")

				if port != 4034 {
					t.Errorf("\t\t%s Should have key %q with value %d", failed, "PORT", 4034)
				} else {
					t.Logf("\t\t%s Should have key %q with value %d", succeed, "PORT", 4034)
				}
			}

			flag, err := cfg.Bool("FLAG")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "FLAG")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", succeed, "FLAG")

				if flag == false {
					t.Errorf("\t\t%s Should have key %q with value %v", failed, "FLAG", true)
				} else {
					t.Logf("\t\t%s Should have key %q with value %v", succeed, "FLAG", true)
				}
			}

			u, err := cfg.URL("DSN")

			if err != nil {
				t.Errorf("\t\t%s Should not return error when valid key %q", failed, "DSN")
			} else {
				t.Logf("\t\t%s Should not return error when valid key %q", succeed, "DSN")

				if u.String() != uStr {
					t.Errorf("\t\t%s Should have key %q with value %v", failed, "DSN", true)
				} else {
					t.Logf("\t\t%s Should have key %q with value %v", succeed, "DSN", true)
				}
			}
		}
	}
}

// TestNotExists validates the ability to load configuration values
// using the OS-level environment variables and panic when something
// is missing.
func TestNotExists(t *testing.T) {
	t.Log("Given the need to panic when environment variables are missing.")
	{
		os.Setenv("MYAPP_PROC_ID", "322")
		os.Setenv("MYAPP_SOCKET", "./tmp/sockets.po")
		os.Setenv("MYAPP_PORT", "4034")
		os.Setenv("MYAPP_FLAG", "true")

		cfg.Init("MYAPP")

		t.Log("\tWhen given a namspace key to search for that does NOT exist.")
		{
			shouldPanic(t, "STAMP", func() {
				cfg.MustTime("STAMP")
			})

			shouldPanic(t, "PID", func() {
				cfg.MustInt("PID")
			})

			shouldPanic(t, "DEST", func() {
				cfg.MustString("DEST")
			})

			shouldPanic(t, "ACTIVE", func() {
				cfg.MustBool("ACTIVE")
			})

			shouldPanic(t, "SOCKET_DSN", func() {
				cfg.MustURL("SOCKET_DSN")
			})
		}
	}
}

// shouldPanic receives a context string and a function to run, if the function
// panics, it is considered a success else a failure.
func shouldPanic(t *testing.T, context string, fx func()) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("\t\t%s Should paniced when giving unknown key %q.", failed, context)
		} else {
			t.Logf("\t\t%s Should paniced when giving unknown key %q.", succeed, context)
		}
	}()

	fx()
}
