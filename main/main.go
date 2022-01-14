package main

import (
	"encoding/json"
	"fmt"
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

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Id             int      `json:"id"`
	DatesLocations []string `json:"datesLocations"`
}

func main() {
	artists, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

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

	relation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	artistsData, err := ioutil.ReadAll(artists.Body)

	if err != nil {
		log.Fatal(err)
	}

	locationsData, err := ioutil.ReadAll(locations.Body)

	if err != nil {
		log.Fatal(err)
	}

	datesData, err := ioutil.ReadAll(dates.Body)

	if err != nil {
		log.Fatal(err)
	}

	relationData, err := ioutil.ReadAll(relation.Body)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(locationsData))

	var artistsObject []Artist
	json.Unmarshal(artistsData, &artistsObject)

	var locationsObject []Location
	json.Unmarshal(locationsData, &locationsObject)

	var datesObject []Date
	json.Unmarshal(datesData, &datesObject)

	var relationObject []Relation
	json.Unmarshal(relationData, &relationObject)

	fmt.Println(locationsObject)

}
