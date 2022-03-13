package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Calendar struct {
	EventID   int
	UserID    int
	Available []struct {
		start string
		end   string
	}
}

var calendars []Calendar

func RemoveCalendar(s []Calendar, index int) []Calendar {
	ret := make([]Calendar, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}

func retrieveCalendar(eventID int, userID int) (calendar Calendar, err error) {
	for _, calendar := range calendars {
		if calendar.EventID == eventID && calendar.UserID == userID {
			return calendar, nil
		}
	}

	return Calendar{0, 0, nil}, errors.New("User not found")
}

func removeCalendar(eventID int, userID int) {
	for index, calendar := range calendars {
		if calendar.EventID == eventID && calendar.UserID == userID {
			calendars = RemoveCalendar(calendars, index)
    		return
		}
	}
}

func getCalendarHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	eventID, _ := strconv.Atoi(r.Form.Get("EventID"))
	userID, _ := strconv.Atoi((r.Form.Get("UserID"))
	
	calendar, err := retrieveCalendar(eventID, userID)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	calendarBytes, err := json.Marshal(calendar)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Write(calendarBytes)
}

func setCalendarHandler(w http.ResponseWriter, r *http.Request) {
	calendar := Calendar{}

	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	calendar.EventID = r.Form.Get("EventID")
	calendar.UserID = r.Form.Get("UserID")
	calendar.Available = r.Form.Get("Available")

	removeCalendar(eventID, userID)
	calendars = append(calendars, calendar)

	w.WriteHeader(http.StatusFound)
}
