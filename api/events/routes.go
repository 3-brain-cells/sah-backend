package events

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/3-brain-cells/sah-backend/util"
	"github.com/go-chi/chi"
)

func Routes(database db.Provider) *chi.Mux {
	router := chi.NewRouter()
	router.Put("/", CreatePartialEvent(database))
	router.Put("/{id}", PopulateEvent(database))
	router.Get("/{id}/vote_options", GetVoteOptions(database))
	router.Post("/{id}/votes", PostVotes(database))

	return router
}

// when users use the slash command (but have not yet populated other fields, times, etc.)
func CreatePartialEvent(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var partialEvent types.EventCreate
		err := json.NewDecoder(r.Body).Decode(&partialEvent)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		err = eventProvider.CreatePartial(r.Context(), partialEvent)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
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

// need to confirm that the user who is populating the event is the same as the user who created the event
func PopulateEvent(eventProvider db.EventProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")
		if id == "" {
			util.ErrorWithCode(r, w, errors.New("the URL parameter is empty"),
				http.StatusBadRequest)
			return
		}

		var event types.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		err = eventProvider.PopulateEvent(r.Context(), event)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
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

		var votes types.UserVotes
		err := json.NewDecoder(r.Body).Decode(&votes)
		if err != nil {
			util.ErrorWithCode(r, w, err, http.StatusBadRequest)
			return
		}

		err = eventProvider.PostVotes(r.Context(), votes, id)
		if err != nil {
			util.Error(r, w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
