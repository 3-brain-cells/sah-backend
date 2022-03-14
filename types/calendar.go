package types

type Time struct {
	start string `json:"start"`
	end   string `json:"end"`
}

type CalendarResponse struct {
	EventID   int    `json:"eventID"`
	UserID    int    `json:"userID"`
	Available []Time `json:"available"`
}
