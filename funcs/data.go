package funcs

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

var all_artists []Artist

var searchResult []SearchRes

// An http client that will be used to get the json data from the API
var client *http.Client

// Template variable used to execute the html files with their output
var tpl *template.Template

// Struct to hold the status code and msg in case of any error
type Result struct {
	Code   int
	Status string
}

// Used to hold multiple instances of Artist interface
type ArtistArray struct {
	Artists   []Artist
	Locations []string
	Empty     bool
}

// Holds all the information for each artist from the API including other APIs for locations,dates, and relations
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
	LocArray     []string
}

// To hold the data from the relations API
type Relation struct {
	ID             int `json:"id"`
	DatesLocations map[string][]string
}

// To hold the data from the dates API
type Dates struct {
	ID    int
	Dates []string
}

// To hold the data from the dates API
type Location struct {
	ID        int
	Locations []string
	Dates     string
}

// Holds the combined result from all the APIs data after processing. Data to be displayed
type Output struct {
	Details Artist
	Rel     Relation
	Loc     Location
	Date    Dates
}

type SearchRes struct {
	ID       int
	Artist   string
	Location string
	Member   string
}

type Coordinates struct {
	Loc string
	Lat float64
	Lng float64
}

type GeocodingResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

type MapResult struct {
	ID int
	LocCoords []Coordinates
}

// Sets the client and template to be used inside the package
func SetTC(t *template.Template, c *http.Client) {
	client = c
	tpl = t
}

// Fetches artists data from the API
func FetchArtists() string {
	artists, err := GetArtists()
	if err != "" {
		return err
	}
	all_artists = artists
	for i, c := range artists {
		loc, err1 := GetLocations(c.Locations)
		if err1 != "" {
			log.Fatal(err1)
		}
		artists[i].LocArray = loc.Locations
	}
	return ""
}

/*
Below functions are used to fetch the data from the given APIs.
Each function will save the fetched data into the right instance of the struct.
They call the GetJson dunction that receives a url and gets the json data from it.
*/
func GetArtists() ([]Artist, string) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var artists []Artist
	err := GetJson(url, &artists)
	if err != nil {
		return nil, "Internal server error"
	}
	return artists, ""
}

func GetDates(date string) (Dates, string) {
	url := date
	dates := Dates{}
	err := GetJson(url, &dates)
	if err != nil {
		return dates, "Internal server error"
	}
	return dates, ""
}

func GetRelation(rel string) (Relation, string) {
	url := rel
	relation := Relation{}
	err := GetJson(url, &relation)
	if err != nil {
		return relation, "Internal server error"
	}
	return relation, ""
}

func GetLocations(loc string) (Location, string) {
	url := loc
	places := Location{}
	err := GetJson(url, &places)
	if err != nil {
		return places, "Internal server error"
	}
	return places, ""
}

func GetJson(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
