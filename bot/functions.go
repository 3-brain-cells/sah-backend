package bot

import (
	"context"
	"fmt"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/segmentio/ksuid"
)

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
