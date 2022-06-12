package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/zinthose/latlngtz"
	"github.com/zinthose/latlngtz/sources/bing"
	"github.com/zinthose/latlngtz/sources/geonames"
	"github.com/zinthose/latlngtz/sources/google"
	"github.com/zinthose/latlngtz/sources/timezonemapper"
)

func main() {
	// Load .env file if it exists
	godotenv.Load()

	// If an apiKey is omitted, package will attempt to load from the environment
	// We initialize the api sources we want to use in the order we want them to be used.
	// GeoNames is recommended to be used first as it is free but limited.
	timezonemapper.Init("NA") // <-- TimeZoneMapper is a local source and requires no API Key
	geonames.Init("")         // <-- Will default to "demo" API Key if omitted, but VERY limited
	google.Init("")           // <-- If you want to use Google, you need to provide an API Key
	bing.Init("")             // <-- If you want to use Bing, you need to provide an API Key

	// Get timeZoneId from first source that returns a result
	timeZoneId, sourceName, err := latlngtz.Get(37.8267, -122.4233)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got %q from %s", timeZoneId, sourceName)

}
