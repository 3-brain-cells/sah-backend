package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Location struct {
	EventID   int
	UserID    int
	Latitude  float64
	Longitude float64
	Address   string
}

func RemoveLocation(s []Location, index int) []Location {
	ret := make([]Location, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

var locations []Location

func retrieveLocation(eventID int, userID int) (location Location, err error) {
	for _, location := range locations {
		if location.EventID == eventID && location.UserID == userID {
			return location, nil
		}
	}

	return Location{}, errors.New("Location not found")
}

func removeLocation(eventID int, userID int) {
	for index, location := range locations {
		if location.EventID == eventID && location.UserID == userID {
			locations = RemoveLocation(locations, index)
			return
		}
	}
}

func getLocationHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, _ := strconv.Atoi(r.Form.Get("eventID"))
	userID, _ := strconv.Atoi(r.Form.Get("userID"))

	location, err := retrieveLocation(eventID, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	locationBytes, err := json.Marshal(location)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(locationBytes)
}

func setLocationHandler(w http.ResponseWriter, r *http.Request) {
	location := Location{}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	location.EventID, _ = strconv.Atoi(r.Form.Get("eventID"))
	location.UserID, _ = strconv.Atoi(r.Form.Get("userID"))
	location.Latitude, _ = strconv.ParseFloat(r.Form.Get("latitude"), 64)
	location.Longitude, _ = strconv.ParseFloat(r.Form.Get("longitude"), 64)
	location.Address = r.Form.Get("address")

	removeLocation(location.EventID, location.UserID)
	locations = append(locations, location)

	w.WriteHeader(http.StatusFound)
}
