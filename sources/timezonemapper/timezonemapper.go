package timezonemapper

// "github.com/zsefvlol/timezonemapper"

import (
	"fmt"
	"reflect"

	"github.com/zinthose/latlngtz"
	"github.com/zsefvlol/timezonemapper"
)

// Define package scoped variables
var source *latlngtz.Source
var defaults latlngtz.Source

// Initialize the Source with the given API Key. If the API Key is omitted, the API Key is loaded from the environment.
func Init(apiKey string) error {
	var err error

	// Use "NA" apiKey as no API Key is required
	defaults.APIKey = "NA"

	type empty struct{}
	source, err = source.Init(defaults, apiKey, reflect.TypeOf(empty{}).PkgPath(), Get)
	return err
}

// Get the timezone for the given latitude and longitude from current source
func Get(latitude float64, longitude float64) (string, error) {

	timeZoneId := timezonemapper.LatLngToTimezoneString(latitude, longitude)
	if timeZoneId == "" {
		err := fmt.Errorf("TimeZoneMapper API Error: unable to lookup timeZoneId")
		return "", err
	}

	return timeZoneId, nil
}
