package events

import (
	"context"
	"fmt"
	"time"

	"github.com/3-brain-cells/sah-backend/api/locations"
	"github.com/3-brain-cells/sah-backend/bot"
	"github.com/3-brain-cells/sah-backend/db"
	"github.com/bwmarrin/discordgo"
)

// ManageEvent manages an event after it has been populated
func ManageEvent(eventProvider db.EventProvider, discordSession *discordgo.Session, eventID string) {
	currentTime := time.Now()

	// get the event associated with the eventID
	ctx := context.Background()
	event, err := eventProvider.GetSingle(ctx, eventID)
	if err != nil {
		fmt.Println("error getting event: ", err)
		return
	}
	if currentTime.Before(event.SwitchToVotingTime) {
		// event is currently in scheduling phase
		str := fmt.Sprintf("Event %s is now in scheduling phase. Please input your schedule to the following link https://INSERTLINK.com/%s", event.Title, event.EventID)
		bot.SchedulingMessage(discordSession, str, event.ChannelID)
		time.Sleep(event.SwitchToVotingTime.Sub(currentTime))
	}
	if currentTime.Before(event.EarliestDate) {
		// calculate best time and location options and update the database
		availTimes := FindAvailability(*event)
		availLocations, err := locations.GetNearby(*event)
		if err != nil {
			fmt.Println("error getting locations: ", err)
			return
		}
		// update these two to the database
		event.VoteOptions.StartEndPairs = availTimes
		event.VoteOptions.Location = availLocations
		// update the database
		ctx := context.Background()
		err = eventProvider.UpdateVoteOptions(ctx, event.VoteOptions, event.EventID)
		if err != nil {
			fmt.Println("error updating event: ", err)
			return
		}
		str := fmt.Sprintf("Event %s is now in voting phase. Please vote at the following link https://INSERTLINK.com/%s", event.Title, event.EventID)
		bot.SchedulingMessage(discordSession, str, event.ChannelID)
		time.Sleep(event.EarliestDate.Sub(currentTime))
	}
	// TODO: print final message with hangout time and location
	// ^^ calls Varnika's function to get the final
	str := "TODO"
	bot.SchedulingMessage(discordSession, str, event.ChannelID)

}

// upon restart of the application, need to restart all in progress events
// get all events from the database
// for each event, check if it is in progress (compare the last time to current time and is populated)
// if it is in progress, then restart it (call ManageEvent), else remove it
func Restart(eventProvider db.EventProvider, discordSession *discordgo.Session) {
	// get all events
	ctx := context.Background()

	events, err := eventProvider.GetAllEvents(ctx)
	if err != nil {
		fmt.Println("Error creating event: ", err)
		// return "try again"
	}

	// filter events
	for _, event := range events {
		if event.Populated {
			// event is in progress
			// check that it is still before the initial event time
			if time.Now().Before(event.EarliestDate) {
				// restart event
				ManageEvent(eventProvider, discordSession, event.EventID)
			}
		}
	}

}
