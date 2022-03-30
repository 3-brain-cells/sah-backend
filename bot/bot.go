package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/3-brain-cells/sah-backend/env"
	"github.com/bwmarrin/discordgo"
)

// Bot parameters
// var (
// 	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
// 	AppID          = flag.String("appid", "", "Bot app ID")
// 	BotToken       = flag.String("token", "", "Bot access token")
// 	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
// )

// permissions 397284730944

// var s *discordgo.Session
// var GuildID = ""
// var AppID = ""
// var BotToken = ""
// var RemoveCommands = false

// var (
// 	integerOptionMinValue = 1.0

// 	commands = []*discordgo.ApplicationCommand{
// 		// {
// 		// 	Name:        "help",
// 		// 	Description: "Command that displays all available commands and their functions",
// 		// },
// 		{
// 			Name:        "create-event",
// 			Description: "Command to create an event",
// 		},
// 		// {
// 		// 	Name:        "vote",
// 		// 	Description: "Command to start voting for the event",
// 		// 	Options: []*discordgo.ApplicationCommandOption{
// 		// 		{
// 		// 			Type:        discordgo.ApplicationCommandOptionString,
// 		// 			Name:        "event-name",
// 		// 			Description: "Name of the event to begin voting for",
// 		// 			Required:    true,
// 		// 		},
// 		// 	},
// 		// },
// 	}
// 	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 		// "help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		// 		Data: &discordgo.InteractionResponseData{
// 		// 			Content: "Hey there! Commands you can use:\n 1. /help - display this menu \n 2. /create-event - create an event \n 3. /vote - start voting for a particular event.",
// 		// 		},
// 		// 	})
// 		// },
// 		"create-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionResponseChannelMessageWithSource,
// 				Data: &discordgo.InteractionResponseData{
// 					Content: "Creating new event! Event details at peepeepoopoo!",
// 				},
// 			})
// 		},
// 		// "vote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		// 	margs := []interface{}{
// 		// 		// Need to check if event name is valid. If not, error
// 		// 		// Also need to link to website to fetch event details.
// 		// 		i.ApplicationCommandData().Options[0].StringValue(),
// 		// 	}
// 		// 	// msgformat := `Sorry, this event does not exist :(`

// 		// 	msgformat :=
// 		// 		` Voting for event **"%s"** location and time has started: https://5302-128-61-84-107.ngrok.io/demo/voting/1

// 		// 	**Possible times**
// 		// 	- 3/14 8:00 PM - 9:00 PM
// 		// 	- 3/15 7:00 PM - 8:00 PM
// 		// 	- 3/16 5:00 PM - 6:00 PM

// 		// 	**Possible locations**
// 		// 	- **Skate Park** (Owens Field Skate Park - 1351 Jim Hamilton Blvd, Columbia, SC 29205)
// 		// 	- **Beltine Lanes** (Beltine Lanes, 2154 S Beltline Blvd, Columbia, SC 29201)
// 		// 	- **Blossom Buffet** (2515 Sunset Blvd West Columbia, SC 29169)
// 		// 	- **Massage Therapy by Trudie Harris (232 Skyland Dr, Columbia, SC 29210)`

// 		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		// 		Data: &discordgo.InteractionResponseData{
// 		// 			Content: fmt.Sprintf(
// 		// 				msgformat,
// 		// 				margs...,
// 		// 			),
// 		// 		},
// 		// 	})
// 		// },
// 	}
// )

func run_bot() {

	var s *discordgo.Session
	GuildID := ""
	RemoveCommands := false

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-event",
			Description: "Command to create an event",
		},
	}

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Creating new event! Event details at peepeepoopoo!",
				},
			})
		},
	}

	token, _ := env.GetEnv("token", "BOT_TOKEN")
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	log.Println("Exiting init")

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutdowning")
}
