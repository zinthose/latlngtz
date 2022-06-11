package geonames

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/zinthose/latlngtz"
)

// Define package scoped variables
var source *latlngtz.Source
var defaults latlngtz.Source

// Initialize the Source with the given API Key. If the API Key is omitted, the API Key is loaded from the environment.
func Init(apiKey string) error {
	var err error
	type empty struct{}

	// Use "demo" apiKey if no API Key is given
	defaults.APIKey = "demo"

	source, err = source.Init(defaults, apiKey, reflect.TypeOf(empty{}).PkgPath(), Get)
	return err
}

// GeoNames API Result Structure but reduced to only hold used data.
// It is recommended to ise only the minimum viable data as it will minimize
// any JSON parsing errors
type apiResult struct {
	TimeZoneId string `json:"timeZoneId"`
	Status     struct {
		Message string `json:"message"`
		Value   int    `json:"value"`
	} `json:"status"`
}

// GeoNames API Errors lookup table
var apiErrors map[int]error = map[int]error{
	10: ErrAuthenticationException,
	11: ErrRecordNotFound,
	12: ErrOtherError,
	13: ErrDatabaseTimeout,
	14: ErrInvalidParameter,
	15: ErrNoResultsFound,
	18: ErrAPIDailyLimitExceeded,
	19: ErrAPIHourlyLimitExceeded,
	20: ErrAPIWeeklyLimitExceeded,
	21: ErrInvalidInput,
	22: ErrServerOverloaded,
	23: ErrServiceNotImplemented,
}

// Known GeoNames API Errors See: http://www.geonames.org/export/webservice-exception.html
var ErrAuthenticationException = errors.New("authentication exception")        // 10
var ErrRecordNotFound = errors.New("record does not exist")                    // 11
var ErrOtherError = errors.New("other error")                                  // 12
var ErrDatabaseTimeout = errors.New("database timeout")                        // 13
var ErrInvalidParameter = errors.New("invalid parameter")                      // 14
var ErrNoResultsFound = errors.New("no results found")                         // 15
var ErrAPIDailyLimitExceeded = errors.New("daily limit of credits exceeded")   // 18
var ErrAPIHourlyLimitExceeded = errors.New("hourly limit of credits exceeded") // 19
var ErrAPIWeeklyLimitExceeded = errors.New("weekly limit of credits exceeded") // 20
var ErrInvalidInput = errors.New("invalid input")                              // 21
var ErrServerOverloaded = errors.New("server overloaded")                      // 22
var ErrServiceNotImplemented = errors.New("service not implemented")           // 23

// Get the timezone for the given latitude and longitude from current source
func Get(latitude float64, longitude float64) (string, error) {

	var result *apiResult
	url := fmt.Sprintf("http://api.geonames.org/timezoneJSON?lat=%f&lng=%f&username=%s", latitude, longitude, source.APIKey)
	err := latlngtz.GetAPIResult(url, &result)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}

	// Test if the response is valid
	if result.Status.Value != 0 {
		log.Printf("GeoNames API error #%d: %s\n", result.Status.Value, result.Status.Message)
		// Check if the error is known
		if err, ok := apiErrors[result.Status.Value]; ok {
			return "", err
		}
		return "", fmt.Errorf("GeoNames API error #%d: %s", result.Status.Value, result.Status.Message)
	}

	return result.TimeZoneId, nil
}
