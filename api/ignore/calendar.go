package ignore

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Calendar struct {
	EventID   int
	UserID    int
	Available []struct {
		start time.Time
		end   time.Time
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

	return Calendar{}, errors.New("User not found")
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

	eventID, _ := strconv.Atoi(r.Form.Get("eventID"))
	userID, _ := strconv.Atoi(r.Form.Get("userID"))

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

	calendar.EventID, err = strconv.Atoi(r.Form.Get("eventID"))
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	calendar.UserID, err = strconv.Atoi(r.Form.Get("userID"))
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// listOfAvailable := r.Form.Get("available")
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("Error: %v", err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	removeCalendar(calendar.EventID, calendar.UserID)
	calendars = append(calendars, calendar)

	w.WriteHeader(http.StatusFound)
}
