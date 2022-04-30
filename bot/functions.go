package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
	"github.com/segmentio/ksuid"
)

func create_event(userID string, guildID string, channelID string, eventProvider db.EventProvider) string {
	// generate an eventID
	// as a short, 5-character string of random alphanumeric characters
	rawId := ksuid.New().String()
	eventID := strings.ToLower(rawId[len(rawId)-5:])

	event := types.Event{CreatorID: userID, GuildID: guildID, EventID: eventID, ChannelID: channelID}

	ctx := context.Background()

	err := eventProvider.CreatePartial(ctx, event)
	if err != nil {
		fmt.Println("Error creating event: ", err)
		// return "try again"
	}
	return fmt.Sprintf("New event created: <https://super-auto-hangouts.netlify.app/new/%v>", eventID)
}
