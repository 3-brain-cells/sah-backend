package locations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/3-brain-cells/sah-backend/env"
	"github.com/3-brain-cells/sah-backend/types"
)

// type Location struct {
// 	Name    string `json:"name" bson:"name"`
// 	Address string `json:"address" bson:"address"`
// 	Rating  int    `json:"rating" bson:"rating"`
// 	Image   string `json:"image" bson:"image"`
// }

func GetNearby(event types.Event) ([]types.Location, error) {
	var midpoint types.Coordinates
	midpoint.Latitude = 0
	midpoint.Longitude = 0

	api_key, err := env.GetEnv("api_key", "GOOGLE_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	for _, location := range event.UserLocations {
		// get coordinates from the location
		url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%v&key=%v", location.LocationID, api_key)

		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			continue
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var woo map[string]interface{}

		json.Unmarshal(body, &woo)
		results := woo["results"].([]interface{})
		coords := results[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"].(map[string]interface{})
		lat := coords["lat"].(float64)
		lng := coords["lng"].(float64)

		midpoint.Latitude += lat
		midpoint.Longitude += lng

	}

	midpoint.Latitude = midpoint.Latitude / float64(len(event.UserLocations))
	midpoint.Longitude = midpoint.Longitude / float64(len(event.UserLocations))

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%v,%v&radius=1500&type=restaurant&key=%v", midpoint.Latitude, midpoint.Longitude, api_key)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))

	// parse body into list of locations
	var woo map[string]interface{}

	json.Unmarshal(body, &woo)
	results := woo["results"].([]interface{})
	// get the first 10 results from results
	results = results[:10]
	// parse these to get out rating, name, address, and image
	locations := make([]types.Location, len(results))
	for i, result := range results {
		result := result.(map[string]interface{})
		locations[i].Name = result["name"].(string)
		locations[i].Address = result["vicinity"].(string)
		locations[i].Image = result["icon"].(string)
		locations[i].Rating = result["rating"].(int)
	}
	return locations, nil
	// TODO: format the return
}
