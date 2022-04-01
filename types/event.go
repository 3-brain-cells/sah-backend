package types

type EventResponse struct {
	EventID int   `json:"eventID"`
	UserIDs []int `json:"userIDs"`
}

// we need the user who is creating the event
// we need to generate an event ID
// we need to use this event ID for the link
// we need to store the GuildID (only members in the server can access)
type EventCreate struct {
	UserID  int `json:"userID"`
	GuildID int `json:"guildID"`
	EventID int `json:"eventID"`
}

type Event struct {
	EventID int   `json:"eventID"`
	UserIDs []int `json:"userIDs"`
}
