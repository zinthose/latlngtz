package bing

import (
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
	source, err = source.Init(defaults, apiKey, reflect.TypeOf(empty{}).PkgPath(), Get)
	return err
}

// Bing API Result Structure but reduced to only hold used data.
type apiResult struct {
	StatusCode        int    `json:"statusCode"`
	StatusDescription string `json:"statusDescription"`
	ResourceSets      []struct {
		Resources []struct {
			Type     string `json:"__type"`
			TimeZone struct {
				IanaTimeZoneId string `json:"ianaTimeZoneId"`
			} `json:"timeZone"`
		} `json:"resources"`
	} `json:"resourceSets"`
}

// Get the timezone for the given latitude and longitude from current source
func Get(latitude float64, longitude float64) (string, error) {
	var result *apiResult

	url := fmt.Sprintf("https://dev.virtualearth.net/REST/v1/TimeZone/%.6f,%.6f?key=%s", latitude, longitude, source.APIKey)
	err := latlngtz.GetAPIResult(url, &result)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return "", err
	}

	// Test if the response is valid
	if result.StatusCode != 200 {
		err = fmt.Errorf("Bing API Error: %d - %s", result.StatusCode, result.StatusDescription)
		log.Println(err.Error())
		return "", err
	}

	return result.ResourceSets[0].Resources[0].TimeZone.IanaTimeZoneId, nil
}
