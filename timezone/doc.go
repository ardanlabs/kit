// Package timezone provides two ways to retrieve the timezone for any latitude,longitude position.
//
// Google Time Zone API
//
// https://developers.google.com/maps/documentation/timezone/
//
// This API is free and limits you to 2,500 calls a day. You will receive a GoogleTimezone object if the call is successful.
//
// 		// RetrieveGoogleTimezone calls the Google API to retrieve the timezone for the lat/lng
//		latitude: The latitude position of the location
//		longitude: The longitude position of the location
//
//		// RetrieveGoogleTimezone(latitude float64, longitude float64) (googleTimezone *GoogleTimezone, err error)
//		DstOffset:    Offset for daylight-savings time in seconds. This will be zero if the time zone is not in Daylight Savings Time during the specified timestamp.
//		RawOffset:    Offset from UTC (in seconds) for the given location. This does not take into effect daylight savings.
//		TimezoneID:   Contains the ID of the time zone, such as "America/Los_Angeles" or "Australia/Sydney".
//		TimezoneName: Contains the long form name of the time zone. This field will be localized if the language parameter is set. eg. "Pacific Daylight Time" or "Australian Eastern Daylight Time"
//		Status:       Indicates the status of the response.
//			OK               : Indicates that the request was successful.
//			INVALID_REQUEST  : Indicates that the request was malformed.
//			OVER_QUERY_LIMIT : Indicates the requestor has exceeded quota.
//			REQUEST_DENIED   : Indicates that the the API did not complete the request. Confirm that the request was sent over http instead of https.
//			UNKNOWN_ERROR    : Indicates an unknown error.
//			ZERO_RESULTS     : Indicates that no time zone data could be found for the specified position or time. Confirm that the request is for a location on land, and not over water.
//
// GeoNames Time Zone API
//
// http://www.geonames.org/
//
// This API is free with registration and has no limit. If you upgrade to a premium account you can be guaranteed fast and reliable response times.
//
//		// RetrieveGeoNamesTimezone calls the GeoNames API to retrieve the timezone for the lat/lng
//		latitude: The latitude position of the location
//		longitude: The longitude position of the location
//		userName: The GeoNames username for using the API
//
// 		// RetrieveGeoNamesTimezone(latitude float64, longitude float64, userName string) (geoNamesTimezone *GeoNamesTimezone, err error)
// 		Time:        	The local current time.
// 		CountryName: 	ISO 3166 country code name
// 		CountryCode: 	ISO 3166 country code
// 		Sunset:      	Current days time for sunset
// 		RawOffset:   	The amount of time in hours to add to UTC to get standard time in this time zone.
// 		DstOffset:   	Offset to GMT at 1. July (deprecated)
// 		GmtOffset:   	Offset to GMT at 1. January (deprecated)
// 		Sunrise:     	Current days time for sunrise
// 		TimezoneID:  	The name of the timezone (according to olson).
// 		Longitude:   	Longitude used for the call
// 		Latitude:    	Latitude used for the call
//
// Using The Time Package
//
// The TimezoneID is a location name corresponding to the IANA Time Zone database. Go provides this database in
// the /lib/time/zoneinfo.zip file. Because of this we can use the TimezoneID within the time package:
//
// 		// Take the timezone id from the API call and create a location object
// 		var location *time.Location
// 		location, err = time.LoadLocation(TimezoneID)
//
//		// Use the location object to make sure this time is in the correct timezone
// 		localTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, location)
//
// 		// Convert to UTC, if required, based on the correct timezone
// 		utcTime := localTime.UTC()
package timezone
