package types

type LocationResponse struct {
	EventID   int     `json:"eventID"`
	UserID    int     `json:"userID"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
