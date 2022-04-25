package events

import (
	"context"
	"fmt"
	"time"

	"github.com/3-brain-cells/sah-backend/db"
)

// ManageEvent manages an event after it has been populated
func ManageEvent(eventProvider db.EventProvider, eventID string) {
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
		// TODO: print out message to tell users to input schedule
		time.Sleep(event.SwitchToVotingTime.Sub(currentTime))
	}
	if currentTime.Before(event.EarliestDate) {
		// TODO: call Varnika's function to calculate best options
		// TODO: print out message to tell users to vote
		time.Sleep(event.EarliestDate.Sub(currentTime))
	}
	// TODO: print final message with hangout time and location
}

// upon restart of the application, need to restart all in progress events
// get all events from the database
// for each event, check if it is in progress (compare the last time to current time and is populated)
// if it is in progress, then restart it (call ManageEvent), else remove it
func Restart(eventProvider db.EventProvider) {
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
				ManageEvent(eventProvider, event.EventID)
			} else {
				// TODO: remove event if we plan to
				// eventProvider.DeleteEvent(ctx)
			}
		}
	}

}
