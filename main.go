package main

import (
	"math/rand"
	"os"
	"posts/v1/lib"
	"posts/v1/server"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zenazn/goji/graceful"
)

func main() {
	// [Hatchways]
	_ = lib.SetupHatchways()

	// initialize seed
	rand.Seed(time.Now().Unix())

	bind := os.Getenv("BIND")
	environment := os.Getenv("APP_ENV")

	if bind == "" {
		bind = "0.0.0.0:8080"
	}

	if environment == "" {
		environment = "development"
	}

	// new web hander
	h, err := server.New(server.Debug(environment != "production"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start web server")
	}

	log.Info().Msgf("[%s] API starting on %s", environment, bind)
	if err := graceful.ListenAndServe(bind, h.Routes()); err != nil {
		log.Fatal().Err(err).Msg("cannot bind to host")
	}

	graceful.Wait()
}
