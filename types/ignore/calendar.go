package ignore

type Time struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type CalendarResponse struct {
	EventID   int    `json:"eventID"`
	UserID    int    `json:"userID"`
	Available []Time `json:"available"`
}
