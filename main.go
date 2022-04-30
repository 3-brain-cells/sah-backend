// functions required
// set event data
// set location data
// set calendar data
// set voting data
// get voting options
// get voting results
// get location data --> determine closest location
// fetch yelp / google maps data
// alogrithm to determine best locations
package main

import (
	"context"
	"flag"
	stdlog "log"
	"os"
	"time"

	"github.com/3-brain-cells/sah-backend/bot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	envPath := flag.String("env", "", "path to .env file")
	logFormat := flag.String("log-format", "console", "log format (one of 'json', 'console')")
	flag.Parse()

	// Set up structured logging
	zerolog.TimeFieldFormat = time.RFC3339Nano
	var logger zerolog.Logger
	switch *logFormat {
	case "console":
		output := zerolog.ConsoleWriter{Out: os.Stdout}
		logger = zerolog.New(output).With().Timestamp().Logger()
	case "json":
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	default:
		log.Fatal().Str("log_format", *logFormat).Msg("unknown log format given")
	}
	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	// Load the .env file if it is specified
	if envPath != nil && *envPath != "" {
		err := godotenv.Load(*envPath)
		if err != nil {
			logger.Fatal().Err(err).Str("env_path", *envPath).Msg("error loading .env file")
		} else {
			logger.Info().Str("env_path", *envPath).Msg("loaded environment variables from file")
		}
	}

	// print out BOT_TOKEN from env file
	logger.Info().Str("BOT_TOKEN", os.Getenv("BOT_TOKEN")).Msg("BOT_TOKEN")
	logger.Info().Str("MONGO_DB_USERNAME", os.Getenv("MONGO_DB_USERNAME")).Msg("MONGO_DB_USERNAME")

	api, err := NewAPIServer(logger)
	if err != nil {
		stdlog.Panicf("error: %v", err)
	}

	ctx := context.Background()

	err = api.Connect(ctx)
	if err != nil {
		stdlog.Panicf("error: %v", err)
	}

	go api.Serve(ctx, 5000)
	// Set up the bot
	bot.RunBot(api.dbProvider, api.discordSession)
}
