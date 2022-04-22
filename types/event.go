package types

import "time"

// we need the user who is creating the event
// we need to generate an event ID
// we need to use this event ID for the link
// we need to store the GuildID (only members in the server can access)
type Event struct {
	CreatorID string `json:"creator_id" bson:"creator_id"`
	GuildID   string `json:"guild_id" bson:"guild_id"`
	EventID   string `json:"id" bson:"id"`

	Title            string           `json:"title" bson:"title"`
	Description      string           `json:"description" bson:"description"`
	EarliestDate     time.Time        `json:"earliest_date" bson:"earliest_date"` // ISO 8601 string
	LatestDate       time.Time        `json:"latest_date" bson:"latest_date"`     // ISO 8601 string
	StartTimeHour    int              `json:"start_time_hour" bson:"start_time_hour"`
	StartTimeMinute  int              `json:"start_time_minute" bson:"start_time_minute"`
	EndTimeHour      int              `json:"end_time_hour" bson:"end_time_hour"`
	EndTimeMinute    int              `json:"end_time_minute" bson:"end_time_minute"`
	LocationCategory LocationCategory `json:"location_category" bson:"location_category"`

	Populated   bool        `json:"populated" bson:"populated"`       // field is set once creator goes on web and populates
	VoteOptions VoteOption  `json:"vote_options" bson:"vote_options"` // ^ not done until this is done
	UserAvailability []UserAvailability `json:"user_availability" bson:"user_availability"` // ^ not done until this is done
	UserVotes   []UserVotes `json:"user_votes" bson:"user_votes"`     // ^ not done until this is done
}

type VoteOption struct {
	Location      []Location `json:"address" bson:"address"`
	StartEndPairs []TimePair `json:"startEndPairs" bson:"startEndPairs"`
}

type Location struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type TimePair struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

type UserVotes struct {
	UserID        string `json:"userID" bson:"userID"`
	LocationVotes []int  `json:"locationVotes" bson:"locationVotes"`
	TimeVotes     []int  `json:"timeVotes" bson:"timeVotes"`
}

type UserAvailability struct {
	UserID string `json:"user_id" bson:"user_id"`
	DayAvailability []DayAvailability `json:"day_availability" bson:"day_availability"`
}

type DayAvailability struct {
	Date time.Time `json:"date" bson:"date"` // ISO 8601 string
	AvailableBlocks []AvailabilityBlock `json:"available_blocks" bson:"available_blocks"`
}

type AvailabilityBlock struct {
	StartHour int `json:"start_hour" bson:"start_hour"`
	StartMinute int `json:"start_minute" bson:"start_minute"`
	EndHour int `json:"end_hour" bson:"end_hour"`
	EndMinute int `json:"end_minute" bson:"end_minute"`
}

type LocationCategory string

var (
	LocationCategoryGeneral LocationCategory = "general"
)
