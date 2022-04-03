package bot

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/3-brain-cells/sah-backend/env"
	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

// Constraints (make function similar to):
func ExampleRunFunction(ctx context.Context, dbProvider db.Provider) error { return nil }

// permissions 397284730944

// var (
// 	integerOptionMinValue = 1.0

// 	commands = []*discordgo.ApplicationCommand{
// 		{
// 			Name:        "create-event",
// 			Description: "Command to create an event",
// 		},
// 	}
// 	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
// 		"create-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 			/* TODO
// 			call create_event function
// 				- should read the user (i.Interaction.User.Username) that calls the function
// 				- should return the event link
// 				- update the GuildID:UserID:Link to database
// 			*/

// 			userID := i.Interaction.Member.User.ID
// 			guildID := i.Interaction.GuildID

// 			// TODO: fix
// 			oauth_login(conf);

// 			content := create_event(userID, guildID, _)

// 			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 				Type: discordgo.InteractionResponseChannelMessageWithSource,
// 				Data: &discordgo.InteractionResponseData{
// 					Content: content,
// 				},
// 			})
// 		},
// 	}
// )

func RunBot(dbProvider db.Provider) {
	var s *discordgo.Session

	token, err := env.GetEnv("token", "BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	id, err := env.GetEnv("id", "BOT_ID")
	if err != nil {
		log.Fatal(err)
	}

	secret, err := env.GetEnv("secret", "BOT_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	// log.Println("Exiting init")

	conf := oath_config(id, secret)

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-event",
			Description: "Command to create an event",
		},
	}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			/* TODO
			call create_event function
				- should read the user (i.Interaction.User.Username) that calls the function
				- should return the event link
				- update the GuildID:UserID:Link to database
			*/

			userID := i.Interaction.Member.User.ID
			guildID := i.Interaction.GuildID

			oauth_login(conf)

			content := create_event(userID, guildID, dbProvider)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}

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
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
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

	if *RemoveCommands {
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
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutdowning")
}
