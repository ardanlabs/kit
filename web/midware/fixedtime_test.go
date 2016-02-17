package midware_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ardanlabs/kit/tests"
	"github.com/ardanlabs/kit/web/app"
	"github.com/ardanlabs/kit/web/midware"
)

func init() {
	tests.Init("KIT")
}

// Success and failure markers.
var (
	success = "\u2713"
	failed  = "\u2717"
)

//==============================================================================

// TestFixedTime ensures that a context going through the FixedTime middleware
// gets its Now time set to a fixed point.
func TestFixedTime(t *testing.T) {
	tests.ResetLog()
	defer tests.DisplayLog()

	t.Log("Given the need to set a fixed time on requests we handle")
	{
		expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

		a := app.New(midware.FixedTime(expected))

		var actual time.Time
		a.Handle("GET", "/", func(c *app.Context) error {
			actual = c.Now
			return nil
		})

		w := httptest.NewRecorder()
		r := tests.NewRequest("GET", "/", nil)

		a.ServeHTTP(w, r)

		t.Log("\tWhen a request is made using the middleware.")
		{
			if actual.Equal(expected) {
				t.Logf("\t\t%s Context.Now should be set to the expected time %v", success, actual)
			} else {
				t.Errorf("\t\t%s Context.Now should be set to the expected time %v", failed, actual)
			}
		}
	}
}
