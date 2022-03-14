package types

type LocationResponse struct {
	EventID   int     `json:"eventID"`
	UserID    int     `json:"userID"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Address   string  `json:"address"`
}
