package types

import "time"

// we need the user who is creating the event
// we need to generate an event ID
// we need to use this event ID for the link
// we need to store the GuildID (only members in the server can access)
type EventCreate struct {
	CreatorID string `json:"creatorID"`
	GuildID   string `json:"guildID"`
	EventID   string `json:"eventID"`
}

type Event struct {
	CreatorID   string      `json:"creatorID"`
	GuildID     string      `json:"guildID"`
	EventID     string      `json:"eventID"`
	VoteOptions VoteOption  `json:"voteOption"`
	UserVotes   []UserVotes `json:"userVotes"`
}

type VoteOption struct {
	Location      []string   `json:"address"`
	StartEndPairs []TimePair `json:"startEndPairs"`
}

type Location struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type TimePair struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type UserVotes struct {
	UserID string     `json:"userID"`
	Votes  []TimePair `json:"votes"`
}
