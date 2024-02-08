package funcs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Handles the default route: HomePage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//Ensure correct http method is used
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" { //Only serve the "/" route, handle 404 error
			res := Result{}
			res.Code = 404
			res.Status = "Not Found"
			ErrorHandler(w, r, &res)
			return
		}

		//Create an instance of the artists array struct to hold all artists
		artists := ArtistArray{}

		artists.Artists = all_artists  //Saves all the artists locally

		//Incase there is no error in fetching the data, execute the index.html file with the artists array
		err := tpl.ExecuteTemplate(w, "index.html", &artists)
		if err != nil {
			log.Fatal(err) //handle any error
		}
	} else if r.Method == "POST" {
		//If the method is http POST method: The form has been submitted to filter the results
		if r.URL.Path != "/"{
			res := Result{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w,r,&res)
			return
		}
		//Call filter function
		Filter(all_artists, r, w)
	} else if r.Method == "SEARCH" {
		//If the method is SEARCH: Requests are being sent using ajax
		//Handles search bar requests using ajax and go routines
		if r.URL.Path != "/"{
			res := Result{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w,r,&res)
			return
		}

		//Will be used to manage functions execution
		var wg sync.WaitGroup
		done := make(chan struct{})
		wg.Add(5)

		query := r.URL.Query().Get("query")  //Retrieve the query value from the url

		searchResult = []SearchRes{}  //Will store all matching objects
		if query != "" {
			go SearchNames(query, &wg)
			go SearchMem(query, &wg)
			go SearchLoc(query, &wg)
			go SearchFirstAlbum(query, &wg)
			go SearchCreationDate(query, &wg)
		}

		//wait utill all functions are completed
		go func() {
			wg.Wait()
			close(done)
		}()
		
		//Once all functions are finished set header and send json response
		<-done
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResult)
	} else {
		res := Result{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w,r,&res)
		return
	}
}

// Handles the "/viewartist" route. Displays the data for the selected artist.
func ViewArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { //Check http method
		if r.URL.Path != "/viewartist" { //Handle path error
			res := Result{}
			res.Code = 404
			res.Status = "Not Found"
			ErrorHandler(w, r, &res)
			return
		}

		//Get the artist id that was passed in the URL with the http GET method
		getID := r.URL.Query().Get("id")

		//Change the string id into int and handle conversion error
		id, err := strconv.Atoi(getID)
		if err != nil {
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

		//Get the list of artists to find the artist with the specified id
		artists := all_artists
		artist := Artist{}
		found := false //Used to indicate whether there exists an artist with the passed id
		for _, c := range artists {
			if c.ID == id {
				found = true
				artist = c
			}
		}

		//If no artist is found, display bad request
		if !found {
			res := Result{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w, r, &res)
			return
		}

		//Fetch the data from the other APIs
		rel, relError := GetRelation(artist.Relations)
		// loc, locError := GetLocations(artist.Locations)
		// dates, dateError := GetDates(artist.ConcertDates)

		//handle fetching errors
		if relError != ""{
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

		modified := make(map[string][]string)
		for key, value := range rel.DatesLocations {
			modified[strings.ToUpper(key)] = value
		}
		rel.DatesLocations = modified

		//Used to correctly format the dates to be displayed
		// for i, c := range dates.Dates {
		// 	dates.Dates[i] = strings.Replace(c, "*", "", -1)
		// }

		if artist.ID == 21 {
			artist.Image = "https://groupietrackers.herokuapp.com/api/images/scorpions.jpeg"
		}

		//Creates the output interface and populates it with the artist's data
		output := Output{}
		output.Details = artist
		output.Rel = rel
		// output.Loc = loc
		// output.Date = dates

		//Executes the viewArtist.html file with the correct output
		err = tpl.ExecuteTemplate(w, "viewArtist.html", &output)
		if err != nil { //Handles execution errors
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	} else {
		//In case the wrong http method was used
		res := Result{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w, r, &res)
		return
	}
}


// Handles the "/map" route. Displays the map page with the locations of the artists
func MapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/map" {
			res := Result{}
			res.Code = 404
			res.Status = "Not Found"
			ErrorHandler(w, r, &res)
			return
		}

		//Get the artist id that was passed in the URL with the http GET method
		getID := r.URL.Query().Get("id")

		//Change the string id into int and handle conversion error
		id, err := strconv.Atoi(getID)
		if err != nil {
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

		//Get the list of artists to find the artist with the specified id
		artists := all_artists
		found := false //Used to indicate whether there exists an artist with the passed id
		for _, c := range artists {
			if c.ID == id {
				found = true
			}
		}

		//If no artist is found, display bad request
		if !found {
			res := Result{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w, r, &res)
			return
		}

		errTpl := tpl.ExecuteTemplate(w, "map.html",nil)
		if errTpl != nil {
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}
	} else {
		res := Result{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w,r,&res)
		return
	}
}


//Handles any requests on the /showmap route
//Calls the Geocode function to get the coordinates of the locations
//Returns a json response with the coordinates
func ShowMap(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		if r.URL.Path != "/showmap" {
			res := Result{}
			res.Code = 404
			res.Status = "Not Found"	
			ErrorHandler(w, r, &res)
			return
		}

		//Get the artist id that was passed in the URL with the http GET method
		getID := r.URL.Query().Get("id")

		//Change the string id into int and handle conversion error
		id, err := strconv.Atoi(getID)
		if err != nil {
			res := Result{}
			res.Code = 500
			res.Status = "Internal Server Error"
			ErrorHandler(w, r, &res)
			return
		}

		//Get the list of artists to find the artist with the specified id
		artists := all_artists
		artist := Artist{}
		found := false //Used to indicate whether there exists an artist with the passed id
		for _, c := range artists {
			if c.ID == id {
				found = true
				artist = c
			}
		}

		//If no artist is found, display bad request
		if !found {
			res := Result{}
			res.Code = 400
			res.Status = "Bad Request"
			ErrorHandler(w, r, &res)
			return
		}

		var locationCoords []Coordinates
		for _, c := range artist.LocArray {
			coords, errGeo := Geocode(c)
			if errGeo != "" {
				res := Result{}
				res.Code = 500
				res.Status = "Internal Server Error"
				ErrorHandler(w, r, &res)
				return
			}
			locationCoords = append(locationCoords, coords)
		}

		res := MapResult{artist.ID, locationCoords}
		w.Header().Set("Content-Type", "application/json")

		// Marshal the MapResult struct to JSON
		mapResultJSON, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Error marshalling to JSON")
			res := Result{Code: 500, Status: "Internal Server Error"}
			ErrorHandler(w, r, &res)
			return
		}
		// Write the JSON response
		w.Write(mapResultJSON)

	} else {
		res := Result{}
		res.Code = 400
		res.Status = "Bad Request"
		ErrorHandler(w,r,&res)
		return
	}
}

// Function to handle any errors. Receives a code and status and executes the error.html file with the corresponding status and status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, res *Result) {
	w.WriteHeader(res.Code)
	err := tpl.ExecuteTemplate(w, "error.html", res)
	if err != nil {
		fmt.Println("Error with error.html")
		os.Exit(2)
	}
}
