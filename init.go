package main

import (
	"digitales-filmmanagement-backend/globals"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

// This function configures the zerolog logging library which is used for
// logging during the initialization and the server start-up. The http
// handlers use another logging library which is more integrated with the
// server.
func init() {
	// set up the time format used in the logging outputs
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// set up error marshalling to show errors with a stacktrace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now read the environment for a variable to set the logging level
	level, levelSet := os.LookupEnv("LOG_LEVEL")
	// now try to parse the level into a zerolog level
	parsedLevel, err := zerolog.ParseLevel(level)

	// now set the level to info if no level was supplied externally,
	// else set the level according to the environment variable if possible.
	// if not, use the default info level
	if !levelSet || err != nil {
		if err != nil {
			log.Warn().Str("readLevel", level).Msg("invalid level supplied in `LOG_LEVEL`. defaulting to 'info'")
		}
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		globals.HttpLogger = globals.HttpLogger.Level(zerolog.InfoLevel)
	} else {
		log.Info().Str("readLevel", level).Msg("configuring zerolog with desired logging level")
		zerolog.SetGlobalLevel(parsedLevel)
		globals.HttpLogger = globals.HttpLogger.Level(parsedLevel)
	}
}
