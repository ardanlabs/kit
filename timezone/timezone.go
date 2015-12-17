package timezone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	googleURI   string = "https://maps.googleapis.com/maps/api/timezone/json?location=%f,%f&timestamp=%d&sensor=false"
	geonamesURI string = "http://api.geonames.org/timezoneJSON?lat=%f&lng=%f&username=%s"
)

//** TYPES

type (
	// GoogleTimezone is the repsonse from the Google timezone API.
	GoogleTimezone struct {
		DstOffset    float64 `bson:"dstOffset"`    // Offset for daylight-savings time in seconds. This will be zero if the time zone is not in Daylight Savings Time during the specified timestamp.
		RawOffset    float64 `bson:"rawOffset"`    // Offset from UTC (in seconds) for the given location. This does not take into effect daylight savings.
		Status       string  `bson:"status"`       // Indicates the status of the response.
		TimezoneID   string  `bson:"timeZoneID"`   // Contains the ID of the time zone, such as "America/Los_Angeles" or "Australia/Sydney".
		TimezoneName string  `bson:"timeZoneName"` // Contains the long form name of the time zone. This field will be localized if the language parameter is set. eg. "Pacific Daylight Time" or "Australian Eastern Daylight Time"
	}

	// GeoNamesTimezone is the repsonse from the GeoNames timezone API.
	GeoNamesTimezone struct {
		Time        string  `bson:"time"`        // The local current time.
		CountryName string  `bson:"countryName"` // ISO 3166 country code name.
		CountryCode string  `bson:"countryCode"` // ISO 3166 country code.
		Sunset      string  `bson:"sunset"`      // Current days time for sunset.
		RawOffset   float64 `bson:"rawOffset"`   // The amount of time in hours to add to UTC to get standard time in this time zone.
		DstOffset   float64 `bson:"dstOffset"`   // Offset to GMT at 1. July (deprecated).
		GmtOffset   float64 `bson:"gmtOffset"`   // Offset to GMT at 1. January (deprecated).
		Sunrise     string  `bson:"sunrise"`     // Current days time for sunrise.
		TimezoneID  string  `bson:"timezoneID"`  // The name of the timezone (according to olson).
		Longitude   float64 `bson:"lng"`         // Longitude used for the call.
		Latitude    float64 `bson:"lat"`         // Latitude used for the call.
	}
)

// RetrieveGoogleTimezone calls the Google API to retrieve the timezone for the lat/lng.
func RetrieveGoogleTimezone(latitude float64, longitude float64) (googleTimezone *GoogleTimezone, err error) {
	defer catchPanic(&err)

	uri := fmt.Sprintf(googleURI, latitude, longitude, time.Now().UTC().Unix())

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	rawDocument, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(rawDocument, &googleTimezone); err != nil {
		return nil, err
	}

	if googleTimezone.Status != "OK" {
		return nil, fmt.Errorf("Error : Google Status : %s", googleTimezone.Status)
	}

	if len(googleTimezone.TimezoneID) == 0 {
		return nil, fmt.Errorf("Error : No Timezone ID Provided")
	}

	return googleTimezone, err
}

// RetrieveGeoNamesTimezone calls the GeoNames API to retrieve the timezone for the lat/lng.
func RetrieveGeoNamesTimezone(latitude float64, longitude float64, userName string) (geoNamesTimezone *GeoNamesTimezone, err error) {
	defer catchPanic(&err)

	uri := fmt.Sprintf(geonamesURI, latitude, longitude, userName)

	resp, err := http.Get(uri)
	if err != nil {
		return geoNamesTimezone, err
	}

	defer resp.Body.Close()

	rawDocument, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return geoNamesTimezone, err
	}

	if err = json.Unmarshal(rawDocument, &geoNamesTimezone); err != nil {
		return geoNamesTimezone, err
	}

	if len(geoNamesTimezone.TimezoneID) == 0 {
		return geoNamesTimezone, fmt.Errorf("Error : No Timezone ID Provided")
	}

	return geoNamesTimezone, err
}

// CatchPanic is used to catch any Panic.
func catchPanic(err *error) {
	if r := recover(); r != nil {
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}
