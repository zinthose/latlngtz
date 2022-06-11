package latlngtz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Source struct {
	Name     string  `json:"name"`
	APIKey   string  `json:"api_key"`
	EnvKey   string  `json:"env_key"`
	Defaults *Source `json:"defaults"`
	Get      SourceGet
}

var Sources map[string]Source

// Initialize the sources
func init() {
	Sources = make(map[string]Source)
}

type SourceGet func(latitude float64, longitude float64) (string, error)

// Load the defaults if not set for this source
func (thisSource *Source) LoadDefaults() {
	if thisSource.Name == "" {
		thisSource.Name = thisSource.Defaults.Name
	}
	if thisSource.APIKey == "" {
		thisSource.APIKey = thisSource.Defaults.APIKey
	}
	if thisSource.EnvKey == "" {
		thisSource.EnvKey = thisSource.Defaults.EnvKey
	}
	thisSource.Get = thisSource.Defaults.Get

}

// Initialize the source
func (thisSource *Source) Init(defaults Source, apiKey string, pkgPath string, getFunc SourceGet) (*Source, error) {
	// Use the package path to get the source name
	pkgPaths := strings.Split(pkgPath, "/")
	pkgName := pkgPaths[len(pkgPaths)-1]
	// Set Defaults
	defaults.Name = pkgName
	defaults.EnvKey = fmt.Sprintf("LATLNGTZ_%s_API_KEY", strings.ToUpper(pkgName))
	defaults.Get = getFunc

	// Load the defaults
	thisSource = &Source{Defaults: &defaults}
	thisSource.LoadDefaults()
	if apiKey == "" {
		if val, ok := os.LookupEnv(thisSource.EnvKey); ok {
			thisSource.APIKey = val
		}
	} else {
		thisSource.APIKey = apiKey
	}

	// If there is no apiKey, ERROR!
	if thisSource.APIKey == "" {
		err := fmt.Errorf("no API key defined for %s API source, please set the API key's environment variable %q", thisSource.Name, thisSource.EnvKey)
		log.Println(err.Error())
		return nil, err
	}

	// Add the source to the sources map
	Sources[thisSource.Name] = *thisSource

	return thisSource, nil
}

// Get the timezone for the given latitude and longitude from first source that returns a timezone
func Get(latitude float64, longitude float64) (string, string, error) {
	if len(Sources) == 0 {
		return "", "", fmt.Errorf("no sources defined, please Init() at least one source or verify the API keys are set")
	}
	for _, source := range Sources {
		timeZoneId, err := source.Get(latitude, longitude)
		if err != nil {
			return "", source.Name, err
		}
		return timeZoneId, source.Name, nil
	}
	return "", "", fmt.Errorf("no sources returned a timezone")
}

// GetAPIResult returns the body of the HTTP response from the given URL
func GetAPIResult(url string, apiResult any) error {

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal result
	err = json.Unmarshal(body, &apiResult)
	if err != nil {
		return err
	}

	// No Errors
	return nil
}
