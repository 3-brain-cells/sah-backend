package events

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

/*
1. user opens link sent by the bot that opens SAH frontend
2. the SAH frontend's javascript runs and finds out it doesnt know who the user is
3. the SAH frontend redirects the user's tab to backend_url/login?event_id=xxxx
4. the SAH backend redirects the user's tab to the oauth flow (discord's server), with some state field set
5. .... something happens in discord
6. Discord redirects the user's tab to SAH backend /auth/callback, with the same state that was set in #4
7. the SAH backend redirects the user's tab to the SAH frontend, and tells it who the user is
*/

// This is the state key used for security, sent in login, validated in callback.
// For this example we keep it simple and hardcode a string
// but in real apps you must provide a proper function that generates a state.
var state = "random"

func oath_config(id string, secret string, router *chi.Mux) {
	// Create a config.
	conf := oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/callback",
		ClientID:     id,
		ClientSecret: secret,
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}

	// step 2: the SAH backend redirects the user's tab to the oauth flow (discord's server),
	// with some state field set
	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		eventID := r.URL.Query()["event_id"]
		http.Redirect(w, r, conf.AuthCodeURL(eventID[0]), http.StatusTemporaryRedirect)
	})

	// Step 6 Discord redirects the user's tab to SAH backend /auth/callback, with the same state that was set in #4
	router.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")

		// Step 3: We exchange the code we got for an access token
		// Then we can use the access token to do actions, limited to scopes we requested
		token, err := conf.Exchange(context.Background(), r.FormValue("code"))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Use the access token, here we use it to get the logged in user's info.
		res, err := conf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

		if err != nil || res.StatusCode != 200 {
			w.WriteHeader(http.StatusInternalServerError)
			if err != nil {
				w.Write([]byte(err.Error()))
			} else {
				w.Write([]byte(res.Status))
			}
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		decodedJson := make(map[string]interface{})
		err = json.Unmarshal(body, decodedJson)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Step 7: the SAH backend redirects the user's tab to the SAH frontend
		// and tells it who the user is
		userID := decodedJson["data"].(map[string]interface{})["id"]
		http.Redirect(w, r, fmt.Sprint("https://super-auto-hangouts.netlify.app/new/%s?user_id=%s", state, userID), http.StatusTemporaryRedirect)
	})

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
