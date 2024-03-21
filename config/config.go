package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnv() {
	if _, exists := os.LookupEnv("DBHOST"); !exists {
		// Attempt to load the .env file if the variable isn't set
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Error loading .env file")
		}
	}
}
