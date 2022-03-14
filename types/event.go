package types

type EventResponse struct {
	EventID int   `json:"eventID"`
	UserIDs []int `json:"userIDs"`
}
