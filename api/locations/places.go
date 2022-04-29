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

	return nil, err

	// TODO: format the return
}
