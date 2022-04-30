package bot

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/3-brain-cells/sah-backend/db"
	"github.com/bwmarrin/discordgo"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

// Constraints (make function similar to):
func ExampleRunFunction(ctx context.Context, dbProvider db.Provider) error { return nil }

func RunBot(dbProvider db.Provider, discordSession *discordgo.Session) {
	// var s *discordgo.Session

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-event",
			Description: "Command to create an event",
		},
	}
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"create-event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			userID := i.Interaction.Member.User.ID
			guildID := i.Interaction.GuildID
			channelID := i.Interaction.ChannelID

			content := fmt.Sprintf("<%s>", create_event(userID, guildID, channelID, dbProvider))

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}

	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	discordSession.AddHandler(func(discordSession *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", discordSession.State.User.Username, discordSession.State.User.Discriminator)
	})

	err := discordSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discordSession.ApplicationCommandCreate(discordSession.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer discordSession.Close()

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
			err := discordSession.ApplicationCommandDelete(discordSession.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutdowning")
}

func SchedulingMessage(discordSession *discordgo.Session, message string, channelID string) {
	_, err := discordSession.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("Cannot send message: %v", err)
	}
}
