package types

import "time"

// we need the user who is creating the event
// we need to generate an event ID
// we need to use this event ID for the link
// we need to store the GuildID (only members in the server can access)
type Event struct {
	CreatorID string `json:"creatorID"`
	GuildID   string `json:"guildID"`
	EventID   string `json:"eventID"`

	Title            string           `json:"title"`
	Description      string           `json:"description"`
	EarliestDate     time.Time        `json:"earliest_date"` // ISO 8601 string
	LatestDate       time.Time        `json:"latest_date"`   // ISO 8601 string
	StartTimeHour    int              `json:"start_time_hour"`
	StartTimeMinute  int              `json:"start_time_minute"`
	EndTimeHour      int              `json:"end_time_hour"`
	EndTimeMinute    int              `json:"end_time_minute"`
	LocationCategory LocationCategory `json:"location_category"`

	Populated   bool        `json:"populated"`  // field is set once creator goes on web and populates
	VoteOptions VoteOption  `json:"voteOption"` // ^ not done until this is done
	UserVotes   []UserVotes `json:"userVotes"`  // ^ not done until this is done
}

type VoteOption struct {
	Location      []Location `json:"address"`
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
	UserID        string `json:"userID"`
	LocationVotes []int  `json:"locationVotes"`
	TimeVotes     []int  `json:"timeVotes"`
}

type LocationCategory string

var (
	LocationCategoryGeneral LocationCategory = "general"
)
