# latlngtz Package
## What is it?
This is a modular system of external API calls to get the timeZoneID (Ex. America/Los_Angeles) from a latitude / longitude coordinate set.

This is also very likely over engineered as I'm learning to program in Go and tend to get carried away with doing things that may not be necessary.

### Currently the following external TimeZone API's are supported:
* [GeoNames](https://www.geonames.org/export/web-services.html#timezone)
* [Google](https://developers.google.com/maps/documentation/timezone/overview)
* [Bing](https://docs.microsoft.com/en-us/bingmaps/rest-services/timezone/find-time-zone)

## Why did I make this?
I'm a field service technician for my day job (I'd rather be programming but meh) and __LOATH__ the current ticketing system being used. So, I'm making my own and this package is just a small piece of the puzzle.

## How to add to your project...
```console
go get github.com/zinthose/latlngtz
````

### Usage

```go
// If an apiKey is omitted, package will attempt to load from the environment
// We initialize the api sources we want to use in the order we want them to be used.
// GeoNames is recommended to be used first as it is free but limited.
geonames.Init("") // <-- Will default to "demo" API Key if omitted, but VERY limited
google.Init("")   // <-- If you want to use Google, you need to provide an API Key
bing.Init("")     // <-- If you want to use Bing, you need to provide an API Key

// Get timeZoneId from first source that returns a result
timeZoneId, sourceName, err := latlngtz.Get(37.8267, -122.4233)
if err != nil {
    panic(err)
}
fmt.Printf("Got %q from %s", timeZoneId, sourceName)
// Result: Got "America/Los_Angeles" from geonames
```

## TODO
 - [ ] Add ability to RoundRobin the API calls to spread the calls across the available sources.
 - [ ] Add Tests 
 - [X] Generate Alpha Release

## Known Issues
* I'm still new to programming in Go (2 weeks of tinkering so far) so expect dragons.  If you see an issue please let me know, or better yet, send a pull request!