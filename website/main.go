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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/3-brain-cells/sah-backend/api"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	startup_context, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	dbConfig := db.DBConfig{
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
	}
	db, err := db.NewDB(startup_context, dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msgf("error init db")
		os.Exit(1)
	}

	serverCtx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Propagate termination signals to the cancellation of the server context
	go func() {
		<-done
		cancel()
	}()

	server := api.NewServer(logger, db)
	server.Serve(serverCtx, 9000)
}

// func index_page_handler(w http.ResponseWriter, r *http.Request) {
// 	// Probably display logo here?
// 	fmt.Fprintf(w, "Super Auto Hangout Backend!")
// }
