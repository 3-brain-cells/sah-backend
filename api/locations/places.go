package locations

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/3-brain-cells/sah-backend/env"
	"github.com/3-brain-cells/sah-backend/types"
)

// type Coordinates struct {
// 	Latitude  float32 `json:"latitude"`
// 	Longitude float32 `json:"longitude"`
// }

func GetNearby(coordinates types.Coordinates) {

	api_key, err := env.GetEnv("api_key", "GOOGLE_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%v,%v&radius=1500&type=restaurant&keyword=cruise&key=%v", coordinates.Latitude, coordinates.Longitude, api_key)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func GetMidpoint(UsersCoordinates []types.Coordinates) types.Coordinates {
	var midpoint types.Coordinates
	midpoint.Latitude = 0
	midpoint.Longitude = 0

	for _, coordinate := range UsersCoordinates {
		midpoint.Latitude += coordinate.Latitude
		midpoint.Longitude += coordinate.Longitude
	}

	midpoint.Latitude = midpoint.Latitude / float32(len(UsersCoordinates))
	midpoint.Longitude = midpoint.Longitude / float32(len(UsersCoordinates))

	return midpoint
}
