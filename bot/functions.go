package bot

import (
	"context"
	"fmt"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/segmentio/ksuid"
)

/* TODO:

NOTE:
we are not doing any location(Google) or Yelp work

1. create_event function
	- take the userID (that is the creator of the event)
	- take the guildID
	- generate an eventID
	- call the function to persist these to the database
	- get a return that this was compelted successfully to the database
	-return the URL for the event page to the server (content in the return of the slash command)

2. discord OAUTH support

3. update messages
	- if the event if populated/created with the web form fields
	- the bot needs to print out the link in the discord and tell people to vote
*/

func create_event(userID string, guildID string, channelID string, eventProvider db.EventProvider) string {

	// generate an eventID
	eventID := ksuid.New().String()

	event := types.Event{CreatorID: userID, GuildID: guildID, EventID: eventID, ChannelID: channelID}

	ctx := context.Background()

	err := eventProvider.CreatePartial(ctx, event)
	if err != nil {
		fmt.Println("Error creating event: ", err)
		// return "try again"
	}
	return fmt.Sprintf("https://super-auto-hangouts.netlify.app/new/%v", eventID)
}
