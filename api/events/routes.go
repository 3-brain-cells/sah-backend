package events

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/3-brain-cells/sah-backend/util"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
)

func Routes(database db.Provider, discordSession *discordgo.Session) *chi.Mux {
	router := chi.NewRouter()

	// create_event ==> CreatePartialEvent() ==>guildID, userID,generate random ID for event ==> put it in to the database
	// user with USERID == creator goes to the sah-hangout.com/{eventID} ==> OAUTH with discord ==> user ID matches ==> fill out the form ==>
	// PUT /{eventID} ==> PopulateEvent() ==> populate the event in the database with all the other stuff
	// GET /{eventID}/voteoptions ==> GetVoteOptions() ==> get the voting options from the database
	// POST /{eventID}/votes ==> PostVotes() ==> OAUTH also ==> post the votes to the database
	// router.Put("/", CreatePartialEvent(database))

	router.Put("/{id}", PopulateEvent(database, discordSession))
	router.Get("/{id}/vote_options", GetVoteOptions(database))
	router.Post("/{id}/votes", PostVotes(database))
	router.Get("/{id}/availability/{user_id}", GetAvailability(database))
	router.Put("/{id}/availability/{user_id}", PutAvailability(database))
	// router.Put("/{id}/location/{user_id}", PutLocation(database))

	return router
}

type GetVoteOptionsResponseBody struct {
	Times     []GetVoteOptionsTime     `json:"times"`
	Locations []GetVoteOptionsLocation `json:"locations"`
}

type GetVoteOptionsTime struct {
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Available []string  `json:"available"`
}

type GetVoteOptionsLocation struct {
	Name                    string  `json:"names"`
	YelpURL                 string  `json:"yelpUrl"`
	Stars                   float64 `json:"stars"`
	DistanceFromCurrentUser float64 `json:"distanceFromCurrentUser"`
	PreviewImageURL         string  `json:"previewImageUrl"`
	Address                 string  `json:"address"`
}

// gets the current events voting options
func GetVoteOptions(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		event, err := eventProvider.GetSingle(r.Context(), id)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		// Convert the data to GetVoteOptionsResponseBody
		responseTimes := make([]GetVoteOptionsTime, len(event.VoteOptions.StartEndPairs))
		for i, time := range event.VoteOptions.StartEndPairs {
			responseTimes[i] = GetVoteOptionsTime{
				Start:     time.Start,
				End:       time.End,
				Available: []string{},
			}
		}
		responseLocations := make([]GetVoteOptionsLocation, len(event.VoteOptions.Location))
		for i, location := range event.VoteOptions.Location {
			responseLocations[i] = GetVoteOptionsLocation{
				Name:                    location.Name,
				YelpURL:                 "",
				Stars:                   4,
				DistanceFromCurrentUser: 8,
				PreviewImageURL:         "",
				Address:                 location.Address,
			}
		}
		responseBody := GetVoteOptionsResponseBody{
			Times:     responseTimes,
			Locations: responseLocations,
		}

		// Return the single announcement as the top-level JSON
		jsonResponse, err := json.Marshal(&responseBody)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

type populateEventRequestBody struct {
	UserID             string                 `json:"user_id"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	EarliestDate       time.Time              `json:"earliest_date"` // ISO 8601 string
	LatestDate         time.Time              `json:"latest_date"`   // ISO 8601 string
	StartTimeHour      int                    `json:"start_time_hour"`
	StartTimeMinute    int                    `json:"start_time_minute"`
	EndTimeHour        int                    `json:"end_time_hour"`
	EndTimeMinute      int                    `json:"end_time_minute"`
	LocationCategory   types.LocationCategory `json:"location_category"`
	SwitchToVotingTime time.Time              `json:"switch_to_voting"` // ISO 8601 string
}

// need to confirm that the user who is populating the event is the same as the user who created the event
func PopulateEvent(eventProvider db.EventProvider, discordSession *discordgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		var body populateEventRequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		// add the field that switches from filling out schedule to voting
		body.SwitchToVotingTime = time.Now().Add(time.Until(body.EarliestDate) / 2)

		// Create the partial event struct.
		// From the provider interface:
		// ignore the following fields:
		// - creatorID
		// - guildID
		// - populated
		// - voteOptions
		// - userVotes
		partialEvent := types.Event{
			EventID:            id,
			Title:              body.Title,
			Description:        body.Description,
			EarliestDate:       body.EarliestDate,
			LatestDate:         body.LatestDate,
			StartTimeHour:      body.StartTimeHour,
			StartTimeMinute:    body.StartTimeMinute,
			EndTimeHour:        body.EndTimeHour,
			EndTimeMinute:      body.EndTimeMinute,
			LocationCategory:   body.LocationCategory,
			SwitchToVotingTime: body.SwitchToVotingTime,
		}

		log.Printf("PopulateEvent event_id=%s user_id=%s", id, body.UserID)
		err = eventProvider.PopulateEvent(r.Context(), partialEvent, body.UserID)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		// create a thread that manages the event
		go ManageEvent(eventProvider, discordSession, partialEvent.EventID)

		w.WriteHeader(http.StatusCreated)
	}

}

type postVotesRequestBody struct {
	UserID        string `json:"user_id"`
	LocationVotes []int  `json:"location_votes"`
	TimeVotes     []int  `json:"time_votes"`
}

// need to confirm that who is putting is in the Guild and has not already voted
func PostVotes(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		var body postVotesRequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		log.Printf("PostVotes event_id=%s user_id=%s", id, body.UserID)
		err = eventProvider.PostVotes(r.Context(), types.UserVotes{
			UserID:        body.UserID,
			LocationVotes: body.LocationVotes,
			TimeVotes:     body.TimeVotes,
		}, id)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

type getAvailabilityResponseBody struct {
	EarliestDate    time.Time `json:"earliest_date"` // ISO 8601 string
	LatestDate      time.Time `json:"latest_date"`   // ISO 8601 string
	StartTimeHour   int       `json:"start_time_hour"`
	StartTimeMinute int       `json:"start_time_minute"`
	EndTimeHour     int       `json:"end_time_hour"`
	EndTimeMinute   int       `json:"end_time_minute"`
	// If null, then availability has not been submitted yet
	Days []types.DayAvailability `json:"days"`
}

func GetAvailability(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the event ID URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			util.ErrorWithCode(r, w, errors.New("the user ID URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		log.Printf("GetAvailability event_id=%s", id)
		event, err := eventProvider.GetSingle(r.Context(), id)
		if err != nil {
			util.Error(r, w, err)
			return
		}
		if event == nil {
			util.ErrorWithCode(r, w, errors.New("event not found"),
				http.StatusNotFound)
			return
		}
		var myAvailabilityDays []types.DayAvailability = nil
		if userAvailability, ok := event.UserAvailability[userID]; ok {
			if len(userAvailability.DayAvailability) > 0 {
				myAvailabilityDays = userAvailability.DayAvailability
			}
		}

		responseBody := getAvailabilityResponseBody{
			EarliestDate:    event.EarliestDate,
			LatestDate:      event.LatestDate,
			StartTimeHour:   event.StartTimeHour,
			StartTimeMinute: event.StartTimeMinute,
			EndTimeHour:     event.EndTimeHour,
			EndTimeMinute:   event.EndTimeMinute,
			Days:            myAvailabilityDays,
		}

		// Return the single announcement as the top-level JSON
		jsonResponse, err := json.Marshal(&responseBody)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

type putLocationRequestBody struct {
	Address string `json:"address"`
}

// func PutLocation(eventProvider db.EventProvider) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := chi.URLParam(r, "id")
// 		if id == "" {
// 			util.ErrorWithCode(r, w, errors.New("the event ID URL parameter is empty"),
// 				http.StatusBadRequest)
// 			return
// 		}

// 		userID := chi.URLParam(r, "user_id")
// 		if userID == "" {
// 			util.ErrorWithCode(r, w, errors.New("the user ID URL parameter is empty"),
// 				http.StatusBadRequest)
// 			return
// 		}

// 		var body putLocationRequestBody
// 		err := json.NewDecoder(r.Body).Decode(&body)
// 		if err != nil {
// 			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
// 			return
// 		}

// 		log.Printf("PutAvailability event_id=%s user_id=%s", id, userID)
// 		err = eventProvider.PutLocation(r.Context(), userID, types.UserLocation{
// 			LocationID: body.Address,
// 		}, id)
// 		if err != nil {
// 			util.Error(r, w, err)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 	}
// }

type putAvailabilityRequestBody struct {
	Days []types.DayAvailability `json:"days"`
	Location types.UserLocation `json:"location"`
}

func PutAvailability(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the event ID URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			util.ErrorWithCode(r, w, errors.New("the user ID URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		var body putAvailabilityRequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		log.Printf("PutAvailability event_id=%s user_id=%s", id, userID)
		err = eventProvider.PutUserAvailabilityAndLocation(r.Context(), userID, types.UserAvailability{
			DayAvailability: body.Days,
		}, body.Location, id)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func FindAvailability(event types.Event) []types.TimePair {
	// 1. divide time slots into buckets of 30 min in the day
	// 2. look for the longest time slots available / most popular
	// 3. first we need to create the buckets and then fill it in
	//    with number of users available

	// 24 * 2 buckets = 48 buckets * numDays
	tempEarliest := time.Date(event.EarliestDate.Year(), event.EarliestDate.Month(), event.EarliestDate.Day(), 0, 0, 0, 0, event.EarliestDate.Location())
	tempLatest := time.Date(event.EarliestDate.Year(), event.LatestDate.Month(), event.LatestDate.Day(), 0, 0, 0, 0, event.LatestDate.Location())
	duration := tempLatest.Sub(tempEarliest)
	var numDays float64 = 1 + duration.Hours()/24
	var buckets = make([]int, int(48*numDays)+1)
	max := 0

	// populate the buckets
	for _, userAvailability := range event.UserAvailability {
		for _, dayAvailability := range userAvailability.DayAvailability {
			// bucket indexing [(dayNum * 48) + (hour * 2)]
			// make temp time, copies event.Earliest Date and makes hour = 0
			offset := dayAvailability.Date.Sub(tempEarliest).Hours() * 2

			for _, block := range dayAvailability.AvailableBlocks {
				startBucket := (block.StartMinute % 30) + (block.StartHour * 2) + int(offset)
				endBucket := int(math.Ceil(float64(block.EndMinute)/30)) + (block.EndHour * 2) + int(offset)
				for i := startBucket; i <= endBucket; i++ {
					buckets[i] += 1
					if buckets[i] > max {
						max = buckets[i]
					}
				}
			}
		}
	}

	startBucket := -1
	var endBucket int
	var ret []types.DayAvailability
	//check which bucket ranges are most popular --> should be == max
	for len(ret) < 3 && max > 0 {
		ret = []types.DayAvailability{}
		for j := 0; j < int(numDays); j++ {
			startBucket = -1
			var dayBlock []types.AvailabilityBlock
			for i := 0; i < 48; i++ {
				if startBucket == -1 && buckets[j*48+i] >= max {
					startBucket = i
				} else if startBucket != -1 && buckets[j*48+i] < max && i-startBucket < 2 {
					// needs to be at least 1 hr long
					continue;
				} else if startBucket != -1 && (buckets[j*48+i] < max || i-startBucket == 12) {
					// max of 3 hrs
					endBucket = i
					var d types.AvailabilityBlock
					d.StartHour = int(startBucket / 2)
					d.EndHour = int(endBucket / 2)
					d.StartMinute = 30 * (startBucket % 2)
					d.EndMinute = 30 * (endBucket % 2)
					dayBlock = append(dayBlock, d)
					startBucket = -1
				}
			}
			var day types.DayAvailability
			day.Date = event.EarliestDate.Add(time.Hour * time.Duration(24*j))

			// check if the dayBlock is empty
			if len(dayBlock) > 0 {
				day.AvailableBlocks = dayBlock
				ret = append(ret, day)
			}
		}
		max = max - 1
	}

	// convert ret to list of time pairs
	var ret2 []types.TimePair
	for _, day := range ret {
		for _, block := range day.AvailableBlocks {
			var pair types.TimePair
			pair.Start = time.Date(day.Date.Year(), day.Date.Month(), day.Date.Day(), block.StartHour, block.StartMinute, 0, 0, day.Date.Location())
			pair.End = time.Date(day.Date.Year(), day.Date.Month(), day.Date.Day(), block.EndHour, block.EndMinute, 0, 0, day.Date.Location())
			ret2 = append(ret2, pair)
		}
	}
	return ret2
}
