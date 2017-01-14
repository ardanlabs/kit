// Package timezone tests the different timezone API calls.
package timezone_test

import (
	"fmt"

	"github.com/ardanlabs/kit/timezone"
)

// ExampleRetrieveGoogleTimezone tests the call to the google timezone API.
func ExampleRetrieveGoogleTimezone() {
	googleTimezone, err := timezone.RetrieveGoogleTimezone(25.7877, -80.2241)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v", googleTimezone)
	// Output:
	// &{0 -18000 OK America/New_York Eastern Standard Time}
}

// ExampleRetrieveGeoNamesTimezone tests the call to the geonames timezone API.
func ExampleRetrieveGeoNamesTimezone() {
	geoNamesTimezone, err := timezone.RetrieveGeoNamesTimezone(25.7877, -80.2241, "ardanstudios")
	if err != nil {
		fmt.Println(err)
		return
	}

	geoNamesTimezone.Time = "2015-12-17 13:13"
	geoNamesTimezone.Sunrise = "2015-12-17 13:13"
	geoNamesTimezone.Sunset = "2015-12-17 13:13"

	fmt.Printf("%v", geoNamesTimezone)
	// Output:
	// &{2015-12-17 13:13 United States US 2015-12-17 13:13 -5 -4 -5 2015-12-17 13:13 America/New_York 0 0}
}
