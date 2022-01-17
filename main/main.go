package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Artist struct {
	Id            int      `json:"id"`
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Members       []string `json:"members"`
	CreationsDate int      `json:"creationDate"`
	FirstAlbum    string   `json:"firstAlbum"`
	Locations     string   `json:"locations"`
	ConcertDates  string   `json:"concertDates"`
	Relations     string   `json:"relations"`
}

var artistsObject []Artist

//can be used to access object within
// type Index struct {
// 	Index []Location `json:"index"`
// }

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

// too many structs needed to unmarshall
// type Relation struct {
// 	Index int `json:"index"`
// }

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArtists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArtists")

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 Status Not Found", 404)
		return
	}
	err = t.Execute(w, artistsObject)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/artists", returnAllArtists)
	fs := http.FileServer(http.Dir("stylesheets/"))
	http.Handle("/stylesheets/",
		http.StripPrefix("/stylesheets/", fs))
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {

	locations, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	dates, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// relation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")

	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }

	locationsData, err := ioutil.ReadAll(locations.Body)

	if err != nil {
		log.Fatal(err)
	}

	datesData, err := ioutil.ReadAll(dates.Body)

	if err != nil {
		log.Fatal(err)
	}

	// relationData, err := ioutil.ReadAll(relation.Body)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(locationsData))

	var locationsObject []Location
	json.Unmarshal(locationsData[9:len(locationsData)-2], &locationsObject)

	var datesObject []Date
	json.Unmarshal(datesData[9:len(datesData)-2], &datesObject)

	// var relationObject Relation
	// json.Unmarshal(relationData, &relationObject)

	artists, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	artistsData, err := ioutil.ReadAll(artists.Body)

	json.Unmarshal(artistsData, &artistsObject)

	if err != nil {
		log.Fatal(err)
	}

	handleRequests()

}
