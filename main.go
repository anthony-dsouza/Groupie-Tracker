package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Artist struct {
	Id             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	LocationsDates Relation
	FirstAlbum     string `json:"firstAlbum"`
	Locations      string `json:"locations"`
	ConcertDates   string `json:"concertDates"`
	Relations      string `json:"relations"`
}

type Relation struct {
	Id             int                 `json:"id"`
	LocationsDates map[string][]string `json:"datesLocations"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

var artistsObject []Artist
var relationObject Relation
var locationsObject Locations
var datesObject Dates

func returnAllArtists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArtists")

	if r.URL.Path != "/" {
		http.Error(w, "404 Status Not Found", 404)
		return
	}

	searchQuery := r.FormValue("q")

	var searchList []Artist

	for _, x := range artistsObject {
		if strings.Contains(x.Name, searchQuery) {
			searchList = append(searchList, x)
		}
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 Status Not Found", 404)
		return
	}
	err = t.Execute(w, searchList)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}

	// http.HandleFunc("/", returnAllArtists)

	if r.FormValue("q") != searchQuery {
		fmt.Println(r.FormValue("q"))
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: searchHandler")

	// u, err := url.Parse(r.URL.String())
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// params := u.Query()
	// searchQuery := params.Get("q")

	searchQuery := r.FormValue("q")

	var searchList []Artist

	for _, x := range artistsObject {
		if strings.Contains(x.Name, searchQuery) {
			searchList = append(searchList, x)
		}
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "404 Status Not Found", 404)
		return
	}
	err = t.Execute(w, searchList)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: profileHandler")

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	profileQuery := params.Get("q")
	idQuery, err := strconv.Atoi(profileQuery)
	if err != nil {
		http.Error(w, "404 Status Not Found - QUERY", 404)
		return
	}

	var profile Artist

	relation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + profileQuery)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	relationData, err := ioutil.ReadAll(relation.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(relationData, &relationObject)

	// assign data to the profile of type Artist using idQuery

	for i, x := range artistsObject {
		if x.Id == idQuery {
			profile = artistsObject[i]
			profile.LocationsDates = relationObject
		}
	}

	t, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "404 Status Not Found - PROFILE", 404)
		return
	}
	err = t.Execute(w, profile)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}

}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: locationsHandler")

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	locationsQuery := params.Get("l")

	locations, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + locationsQuery)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	locationsData, err := ioutil.ReadAll(locations.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(locationsData, &locationsObject)

	t, err := template.ParseFiles("templates/locations.html")
	if err != nil {
		http.Error(w, "404 Status Not Found", 404)
		return
	}
	err = t.Execute(w, locationsObject)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}
}

func datesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: datesHandler")

	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	datesQuery := params.Get("d")

	dates, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + datesQuery)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	datesData, err := ioutil.ReadAll(dates.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(datesData, &datesObject)

	t, err := template.ParseFiles("templates/dates.html")
	if err != nil {
		http.Error(w, "404 Status Not Found", 404)
		return
	}
	err = t.Execute(w, datesObject)
	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)
	}
}

func handleRequests() {
	a := time.Now()
	http.HandleFunc("/", returnAllArtists)
	// http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/dates", datesHandler)
	fs := http.FileServer(http.Dir("stylesheets/"))
	http.Handle("/stylesheets/",
		http.StripPrefix("/stylesheets/", fs))
	fmt.Println("Server Running")
	fmt.Println(time.Since(a))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("500 Internal Server Error")
	}
}

func main() {

	artists, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	artistsData, err := ioutil.ReadAll(artists.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(artistsData, &artistsObject)

	if artistsObject == nil {
		log.Fatal("500 Internal Server Error")
	}

	handleRequests()

}
