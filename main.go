package main

import (
	"fmt"
	"groupie/funcs"
	"html/template"
	"log"
	"net/http"
	"time"
)

// An http client that will be used to get the json data from the API
var client *http.Client

// Template variable used to execute the html files with their output
var tpl *template.Template

// Main function: Will parse the required html files and will handle different routes. Starts a server that listens on port 8080
func main() {
	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// funcs.Geocode()

	//Set a timeout for the client if no response is received within 20 seconds
	client = &http.Client{Timeout: 20 * time.Second}
	funcs.SetTC(tpl, client) //send the client and template to the funcs package
	fmt.Println("Fetching data from API...")
	errf := funcs.FetchArtists()     //Fetch all the artists from the API
	if errf != ""{
		log.Fatal(errf)
	}
	fmt.Println("Data fetched successfully.")
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles")))) //Serving css files
	http.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts")))) //Serving script files

	http.HandleFunc("/", funcs.HomeHandler)       //Handles the default route; the home page
	http.HandleFunc("/viewartist", funcs.ViewArtist)  //Handles the view artist route
	http.HandleFunc("/map", funcs.MapHandler)	 //Handles the map route
	http.HandleFunc("/showmap", funcs.ShowMap)	 //Handles the show map route
	fmt.Println("Starting Server...")
	fmt.Println("Listening on port 8080")
	fmt.Println("Navigate to: http://localhost:8080")
	ServerError := http.ListenAndServe(":8080", nil) //starts a server on port 8080
	if ServerError != nil {
		log.Fatal(ServerError)
	}
}
