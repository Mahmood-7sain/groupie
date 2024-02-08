package funcs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)


//Takes a string location and returns a Coordinates struct and an error message
//Uses googles geocoding api to get the coordinates
func Geocode(location string) (Coordinates, string) {
	var coords Coordinates

	apiKey := "" //Provide your API key to use
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", url.QueryEscape(location), apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return coords, "Error with request"
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return coords, "Error with response"
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return coords, "Error with body"
	}

	var response GeocodingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return coords, "Error with unmarshal"
	}

	if len(response.Results) > 0 {
		coords.Lat = response.Results[0].Geometry.Location.Lat
		coords.Lng = response.Results[0].Geometry.Location.Lng
		coords.Loc = response.Results[0].FormattedAddress
		return coords, ""
	} else {
		return coords, "No results found"
	}
}
