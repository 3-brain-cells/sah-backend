package ignore

type Option struct {
	Time      Time    `json:"time"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Address   string  `json:"address"`
}

type OptionsResponse struct {
	EventID int      `json:"eventID"`
	UserID  int      `json:"userID"`
	Options []Option `json:"options"`
}
