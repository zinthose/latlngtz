package google

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/zinthose/latlngtz"
)

// Define package scoped variables
var source *latlngtz.Source
var defaults latlngtz.Source

// Initialize the Source with the given API Key. If the API Key is omitted, the API Key is loaded from the environment.
func Init(apiKey string) error {
	var err error
	type empty struct{}
	source, err = source.Init(defaults, apiKey, reflect.TypeOf(empty{}).PkgPath(), Get)
	return err
}

// Google API Result Structure but reduced to only hold used data.
// It is recommended to ise only the minimum viable data as it will minimize
// any JSON parsing errors
type apiResult struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	TimeZoneId   string `json:"timeZoneId"`
}

// Get the timezone for the given latitude and longitude from current source
func Get(latitude float64, longitude float64) (string, error) {
	var result *apiResult

	timestamp := time.Now().Unix()
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/timezone/json?location=%.6f,%.6f&timestamp=%d&key=%s", latitude, longitude, timestamp, source.APIKey)
	err := latlngtz.GetAPIResult(url, &result)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}

	// Test if the response is valid
	if result.Status != "OK" {
		err = fmt.Errorf("Google API Error: %s - %s", result.Status, result.ErrorMessage)
		log.Println(err.Error())
		return "", err
	}
	return result.TimeZoneId, nil
}
