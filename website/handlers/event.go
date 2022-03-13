package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Event struct {
	EventID int
	UserIDs []int
}

func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func RemoveEvent(s []Event, index int) []Event {
	ret := make([]Event, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

var events []Event

func retrieveEvent(eventID int) (event Event, err error) {
	for _, event := range events {
		if event.EventID == eventID {
			return event, nil
		}
	}

	return Event{}, errors.New("Event not found")
}

func removeEvent(eventID int) {
	for index, event := range events {
		if event.EventID == eventID {
			events = RemoveEvent(events, index)
			return
		}
	}
}

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, err := strconv.Atoi(r.Form.Get("EventID"))
	userIDsStr := r.Form.Get("UserIDs")

	var userIDs []int

	for _, userIDStr := range userIDsStr {
		userID, _ := strconv.Atoi(userIDStr)
		userIDs = append(userIDs, userID)
	}

	event, err := retrieveEvent(eventID)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(eventBytes)
}

func setEventHandler(w http.ResponseWriter, r *http.Request) {
	event := Event{}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, err := strconv.Atoi(r.Form.Get("EventID"))
	userIDsStr := r.Form.Get("UserIDs")

	var userIDs []int

	for _, userIDStr := range userIDsStr {
		userID, _ := strconv.Atoi(userIDStr)
		userIDs = append(userIDs, userID)
	}

	removeEvent(eventID)
	events = append(events, event)

	w.WriteHeader(http.StatusFound)
}

func addEventMember(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, err := strconv.Atoi(r.Form.Get("EventID"))
	userID, err := strconv.Atoi(r.Form.Get("UserID"))

	for _, event := range events {
		if event.EventID == eventID {
			event.UserIDs = append(event.UserIDs, userID)
		}
	}

	w.WriteHeader(http.StatusFound)
}

func removeEventMember(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, err := strconv.Atoi(r.Form.Get("EventID"))
	userID, err := strconv.Atoi(r.Form.Get("UserID"))

	for _, event := range events {
		if event.EventID == eventID {
			for index, user := range event.UserIDs {
				if user == userID {
					RemoveIndex(event.UserIDs, index)
				}
			}
		}
	}

	w.WriteHeader(http.StatusFound)
}
