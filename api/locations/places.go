package locations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

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
		lat := location.Latitude
		lng := location.Longitude

		fmt.Println("coordinates: %v, %v",lat, lng)


		midpoint.Latitude += lat
		midpoint.Longitude += lng

	}

	midpoint.Latitude = midpoint.Latitude / float64(len(event.UserLocations))
	midpoint.Longitude = midpoint.Longitude / float64(len(event.UserLocations))

	fmt.Println("coordinates: %v, %v", midpoint.Latitude, midpoint.Longitude)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%v,%v&radius=1500&type=restaurant&key=%v", midpoint.Latitude, midpoint.Longitude, api_key)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Dump request
	httputil.DumpRequestOut(req, true)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	// Dump response
	httputil.DumpResponse(res, true)

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
	// parse these to get out rating, name, address, and image
	locations := make([]types.Location, min(len(results), 10))

	for i, result := range results {
		if i < 10 {
			result := result.(map[string]interface{})
			locations[i].Name = result["name"].(string)
			locations[i].Address = result["vicinity"].(string)
			
			if result["photos"] != nil {
				photos := result["photos"].([]interface{})
				if len(photos) > 0 {
					photo := photos[0].(map[string]interface{})["photo_reference"].(string)
					locations[i].Image = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photo_reference=%v&key=%v", photo, api_key)
				}
			} else {
				locations[i].Image = result["icon"].(string)
			}
			locations[i].Rating = results[i].(map[string]interface{})["rating"].(float64)
			locations[i].Latitude = result["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lat"].(float64)
			locations[i].Longitude = result["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lng"].(float64)
		}
	}
	// for i, result := range results {
	// 	result := result.(map[string]interface{})
	// 	locations[i].Name = result["name"].(string)
	// 	locations[i].Address = result["vicinity"].(string)
	// 	locations[i].Image = result["icon"].(string)
	// 	locations[i].Rating = result["rating"].(int)
	// }
	return locations, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}