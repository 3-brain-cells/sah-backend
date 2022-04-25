package events

import (
	"context"
	"fmt"
	"time"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
)

// ManageEvent manages an event after it has been populated
func ManageEvent(partialEvent types.Event) {
	currentTime := time.Now()
	if currentTime.Before(partialEvent.SwitchToVotingTime) {
		// event is currently in scheduling phase
		// TODO: print out message to tell users to input schedule
		time.Sleep(partialEvent.SwitchToVotingTime.Sub(currentTime))
	}
	if currentTime.Before(partialEvent.EarliestDate) {
		// TODO: call Varnika's function to calculate best options
		// TODO: print out message to tell users to vote
		time.Sleep(partialEvent.EarliestDate.Sub(currentTime))
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
				ManageEvent(*event)
			} else {
				// TODO: remove event if we plan to
				// eventProvider.DeleteEvent(ctx)
			}
		}
	}

}
