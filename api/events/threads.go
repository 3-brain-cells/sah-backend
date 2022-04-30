package events

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/3-brain-cells/sah-backend/api/locations"
	"github.com/3-brain-cells/sah-backend/bot"
	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/types"
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
		str := fmt.Sprintf("New event created: **%s**\n"+
			"Possible dates: %v through %v\n"+
			"Possible times: %d:%02d through %d:%02d\n"+
			"\nEnter your availability here: <https://super-auto-hangouts.netlify.app/availability/%s>", event.Title, event.EarliestDate.Format("01-02-2006"), event.LatestDate.Format("01-02-2006"), event.StartTimeHour, event.StartTimeMinute, event.EndTimeHour, event.EndTimeMinute, event.EventID)
		bot.SchedulingMessage(discordSession, str, event.ChannelID)
		// time.Sleep(event.SwitchToVotingTime.Sub(currentTime))
		time.Sleep(time.Second * 60)
	}
	if currentTime.Before(event.EarliestDate) {
		// calculate best time and location options and update the database
		event, err := eventProvider.GetSingle(ctx, eventID)
		if err != nil {
			fmt.Println("error getting event: ", err)
			return
		}
		availTimes := FindAvailability(*event)
		// Add all user colors and names to the vote time options
		addUserColorsAndNames(event.GuildID, availTimes, discordSession)
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
		str := fmt.Sprintf("Voting for event **%s** location and time has started: <https://super-auto-hangouts.netlify.app/vote/%s>\n"+
			"Possible dates: %v through %v\n"+
			"Possible times: %d:%02d through %d:%02d\n", event.Title, event.EventID, event.EarliestDate.Format("01-02-2006"), event.LatestDate.Format("01-02-2006"), event.StartTimeHour, event.StartTimeMinute, event.EndTimeHour, event.EndTimeMinute)
		bot.SchedulingMessage(discordSession, str, event.ChannelID)
		// time.Sleep(event.EarliestDate.Sub(currentTime))
		time.Sleep(time.Second * 60)
	}
	// iterate through all event.uservotes
	// get the location with most votes
	// get the time with most votes

	event, err = eventProvider.GetSingle(ctx, eventID)
	if err != nil {
		fmt.Println("error getting event: ", err)
		return
	}

	if len(event.UserVotes) == 0 {
		log.Printf("No votes for event %s (event_id=%s); returning early", event.Title, event.EventID)
		return
	}

	// make a new uservotes
	var finalTimes []int
	var finalLocations []int
	for _, v := range event.UserVotes {
		finalTimes = make([]int, len(v.TimeVotes))
		finalLocations = make([]int, len(v.LocationVotes))
		break
	}
	for _, userVote := range event.UserVotes {
		// get the location with most votes
		// get the time with most votes
		for i, timeVote := range userVote.TimeVotes {
			finalTimes[i] += timeVote
		}
		for i, locationVote := range userVote.LocationVotes {
			finalLocations[i] += locationVote
		}
	}
	// get the location with most votes
	max := 0
	locationIndex := 0
	for i, locationVote := range finalLocations {
		if locationVote > max {
			max = locationVote
			locationIndex = i
		}
	}

	max = 0
	timeIndex := 0
	for i, timeVote := range finalTimes {
		if timeVote > max {
			max = timeVote
			timeIndex = i
		}
	}

	// get the actual location and time and create string
	locationFinal := event.VoteOptions.Location[locationIndex]
	startEndFinal := event.VoteOptions.StartEndPairs[timeIndex]

	str := fmt.Sprintf("Event %v is now over. The event will take place at %v (%v) from %v till %v", event.Title, locationFinal.Name, locationFinal.Address, startEndFinal.Start, startEndFinal.End)
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

func addUserColorsAndNames(guildID string, availTimes []types.TimePair, discordSession *discordgo.Session) {
	type colorAndName struct {
		Color string
		Name  string
	}

	guild, err := discordSession.Guild(guildID)
	if err != nil {
		fmt.Printf("Error getting guild (guild_id=%s): %v", guildID, err)
	}

	colorMap := make(map[string]colorAndName)
	for i := range availTimes {
		for j := range availTimes[i].Users {
			id := availTimes[i].Users[j].ID
			if _, ok := colorMap[id]; !ok {
				var name string = "unknown"
				var color string = "#222222"

				// Fetch the user's color and name
				member, err := discordSession.GuildMember(guildID, id)
				if err != nil {
					fmt.Printf("Error getting member (user_id=%s, guild_id=%s): %v", id, guildID, err)
				}

				if member != nil {
					if member.Nick != "" {
						name = member.Nick
					} else {
						name = member.User.Username
					}

					colorInt := firstRoleColor(guild, member.Roles)
					if colorInt != 0 {
						color = fmt.Sprintf("#%06X", colorInt)
					}
				}

				// Store the color and name
				colorMap[id] = colorAndName{
					Color: color,
					Name:  name,
				}
			}

			availTimes[i].Users[j].Color = colorMap[id].Color
			availTimes[i].Users[j].Name = colorMap[id].Name
		}
	}
}

// From https://github.com/bwmarrin/discordgo/blob/cd95ccc2d3c030436fcd9ec3caf0b43f539350dd/state.go#L1258
func firstRoleColor(guild *discordgo.Guild, memberRoles []string) int {
	roles := discordgo.Roles(guild.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range memberRoles {
			if role.ID == roleID {
				if role.Color != 0 {
					return role.Color
				}
			}
		}
	}

	for _, role := range roles {
		if role.ID == guild.ID {
			return role.Color
		}
	}

	return 0
}
