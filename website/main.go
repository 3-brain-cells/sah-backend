// functions required
// set event data
// set location data
// set calendar data
// set voting data
// get voting options
// get voting results
// get location data --> determine closest location
// fetch yelp / google maps data
// alogrithm to determine best locations

package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", index_page_handler).Methods("GET")

	// do we need this?
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	// all of them are POST requests because we need to specify the event ID and sometimes, user ID

	// event data
	r.HandleFunc("/get-event", handlers.getEventHandler).Methods("POST")
	r.HandleFunc("/set-event", handlers.setEventHandler).Methods("POST")
	r.HandleFunc("/add-event-member", handlers.addEventMember).Methods("POST")
	r.HandleFunc("/remove-event-member", handlers.setEventHandler).Methods("POST")

	// location data
	r.HandleFunc("/location", handlers.getLocationHandler).Methods("POST")
	r.HandleFunc("/location", handlers.setLocationHandler).Methods("POST")

	// calendar data
	r.HandleFunc("/calendar", handlers.getCalendarHandler).Methods("POST")
	r.HandleFunc("/calendar", handlers.setCalendarHandler).Methods("POST")

	// voting options data
	r.HandleFunc("/voting-options", handlers.getVotingOptionsHandler).Methods("POST")
	r.HandleFunc("/voting-options", handlers.setVotingOptionsHandler).Methods("POST")

	// voting results data
	r.HandleFunc("/voting-results", handlers.getVotingResultsHandler).Methods("POST")
	r.HandleFunc("/voting-results", handlers.setVotingResultsHandler).Methods("POST")

	r.HandleFunc("/best-location", handlers.getBestLocationHandler).Methods("POST")

	// add more?

	return r
}

func main() {
	r := mux.NewRouter()
	http.ListenAndServe(":8080", r)
}

func index_page_handler(w http.ResponseWriter, r *http.Request) {
	// Probably display logo here?
	fmt.Fprintf(w, "Super Auto Hangout Backend!")
}