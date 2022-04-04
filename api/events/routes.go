package events

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/env"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/3-brain-cells/sah-backend/util"
	"github.com/go-chi/chi"
)

func Routes(database db.Provider) *chi.Mux {
	router := chi.NewRouter()

	// create_event ==> CreatePartialEvent() ==>guildID, userID,generate random ID for event ==> put it in to the database
	// user with USERID == creator goes to the sah-hangout.com/{eventID} ==> OAUTH with discord ==> user ID matches ==> fill out the form ==>
	// PUT /{eventID} ==> PopulateEvent() ==> populate the event in the database with all the other stuff
	// GET /{eventID}/voteoptions ==> GetVoteOptions() ==> get the voting options from the database
	// POST /{eventID}/votes ==> PostVotes() ==> OAUTH also ==> post the votes to the database
	// router.Put("/", CreatePartialEvent(database))
	clientID, err := env.GetEnv("token", "BOT_ID")
	if err != nil {
		log.Fatal(err)
	}

	secret, err := env.GetEnv("token", "BOT_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	oath_config(clientID, secret, router)

	router.Put("/{id}", PopulateEvent(database))
	router.Get("/{id}/vote_options", GetVoteOptions(database))
	router.Post("/{id}/votes", PostVotes(database))

	return router
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

		// Return the single announcement as the top-level JSON
		jsonResponse, err := json.Marshal(event.VoteOptions)
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
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	EarliestDate     time.Time              `json:"earliest_date"` // ISO 8601 string
	LatestDate       time.Time              `json:"latest_date"`   // ISO 8601 string
	StartTimeHour    int                    `json:"start_time_hour"`
	StartTimeMinute  int                    `json:"start_time_minute"`
	EndTimeHour      int                    `json:"end_time_hour"`
	EndTimeMinute    int                    `json:"end_time_minute"`
	LocationCategory types.LocationCategory `json:"location_category"`
}

// need to confirm that the user who is populating the event is the same as the user who created the event
func PopulateEvent(eventProvider db.EventProvider) http.HandlerFunc {
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

		// Create the partial event struct.
		// From the provider interface:
		// ignore the following fields:
		// - creatorID
		// - guildID
		// - eventID
		// - populated
		// - voteOptions
		// - userVotes
		partialEvent := types.Event{
			Title:            body.Title,
			Description:      body.Description,
			EarliestDate:     body.EarliestDate,
			LatestDate:       body.LatestDate,
			StartTimeHour:    body.StartTimeHour,
			StartTimeMinute:  body.StartTimeMinute,
			EndTimeHour:      body.EndTimeHour,
			EndTimeMinute:    body.EndTimeMinute,
			LocationCategory: body.LocationCategory,
		}

		err = eventProvider.PopulateEvent(r.Context(), partialEvent /*userID*/, "")
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

type postVotesRequestBody struct {
	LocationVotes []int `json:"location_votes"`
	TimeVotes     []int `json:"time_votes"`
}

// need to confirm that who is putting is in the Guild and has not already voted
func PostVotes(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO authentication, get this user's ID (or maybe don't for demo)
		userID := ""

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

		err = eventProvider.PostVotes(r.Context(), types.UserVotes{
			UserID:        userID,
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
