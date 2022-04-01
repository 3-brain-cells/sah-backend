package routes

import (
	"net/http"
)

type Options struct {
	EventID   int
	PlaceTime []struct {
		time      string
		latitude  float32
		longitude float32
		address   string
	}
}

func getVotingOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// func findLocationsHandler(w http.ResponseWriter, r *http.Request) {
// 	// find the centroid
// 	// 1. convert degrees to radians
// 	// 2. minimize geodesic distance

// 	w.WriteHeader(http.StatusNotImplemented)
// 	urlFormat := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%f,%f&radius=%d&type=%s&keyword=%s&key=%s"
// 	args := []interface{}{
// 		//i.ApplicationCommandData().Options[0].StringValue(),
// 	}
// 	url := fmt.Sprintf(
// 		urlFormat,
// 		args...,
// 	)
// 	method := "GET"

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, nil)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(body))
// }

func getVotingResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func setVotingResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func getBestLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
