package main

import (
	"bet25-calendar-sync/global_state"
	"os"

	"github.com/joho/godotenv"
	log "github.com/s00500/env_logger"
)

func getEnv() *global_state.State {
	// load .env file from given path
	// we keep it empty it will load .env from current directory
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file, %s", err)
	}

	state := global_state.NewState()

	state.DEBUG = os.Getenv("DEBUG")
	state.GOOGLE_CALENDAR_ID = os.Getenv("GOOGLE_CALENDAR_ID")
	state.GOOGLE_CREDENTIALS = os.Getenv("GOOGLE_CREDENTIALS")
	state.GOOGLE_API_KEY = os.Getenv("GOOGLE_API_KEY")

	state.REBOK_USERNAME = os.Getenv("REBOK_USERNAME")
	state.REBOK_PASSWORD = os.Getenv("REBOK_PASSWORD")
	state.REBOK_GET_USER = os.Getenv("REBOK_GET_USER")

	if state.DEBUG == "" {
		state.DEBUG = "false"
	}
	if state.GOOGLE_CALENDAR_ID == "" {
		log.Fatalf("No calendar defined")
	}
	if state.GOOGLE_API_KEY == "" { // Google API Key or credentials is required
		log.Fatalf("No Google API Key defined, checking for Google Credentials")
		if state.GOOGLE_CREDENTIALS == "" {
			log.Fatalf("No Google Credentials defined")
		}
	}

	if state.REBOK_USERNAME == "" {
		log.Fatalf("No Rebook username defined")
	}
	if state.REBOK_PASSWORD == "" {
		log.Fatalf("No Rebook password defined")
	}

	return &state
}
