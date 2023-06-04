package main

import (
	"database/sql"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"digitales-filmmanagement-backend/config"
	"digitales-filmmanagement-backend/globals"
	// database driver
	_ "github.com/go-sql-driver/mysql"
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

// This function reads the configuration file used to set up the database
// connection and the authorization validation
func init() {
	// read the environment and check for the variable `CONFIGURATION_LOCATION`
	fileLocation, isSet := os.LookupEnv("CONFIGURATION_LOCATION")
	// if no location was set manually, use the default value
	if !isSet {
		log.Debug().Str("path", "./configuration.toml").Msg("no configuration location set via environment. using default")
		fileLocation = "./configuration.toml"
	}

	// now create an empty configuration object
	var conf config.Configuration
	// now try to open the configuration file
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open configuration file")
	}
	// and now try to read it
	err = toml.NewDecoder(file).Decode(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to read/parse configuration file")
	}
	// now log the configuration to the debug output
	log.Debug().Interface("config", conf).Send()
	// after reading the configuration, validate the sub-configurations
	err = conf.OIDC.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid OIDC configuration")
	}
	// since this step modifies some values, reprint the configuration
	log.Debug().Interface("config", conf).Msg("configuration modified")
	err = conf.Database.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid database configuration")
	}
	globals.Configuration = conf

	// since the configuration is valid, load the sql queries for the database
	globals.SqlQueries, err = dotsql.LoadFromFile("./queries.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load sql queries")
	}
}

// this function connects the application to the database server and checks if
// the wanted database schema is present on the server. to achieve this, this
// step creates a temporary database connection which connects to a user's
// default database schema. the connection is closed after this step
func init() {
	// get the configuration
	c := globals.Configuration.Database

	// use the configuration to create a dsn without specifying the schema
	dsn := c.BuildSchemalessDSN()

	// now open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open database connection")
	}
	// now defer the closing of the database connection until the function ends
	// if the connection cannot be closed, a warning is printed
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("unable to close database connection")
		} else {
			log.Info().Msg("closed temporary database connection")
		}
	}(db)
	// now log a status message to indicate that the database was connected
	// and the database is now tested
	log.Info().Msg("connected to database via temporary connection")
	log.Info().Msg("checking for database schema")

	// now execute the schema check query
	row, err := globals.SqlQueries.QueryRow(db, "is-schema-available", c.Schema)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to verify the provided database schema name")
	}
	// now parse the response into a boolean
	var schemaFound bool
	err = row.Scan(&schemaFound)

	if !schemaFound {
		log.Fatal().Msg("configured schema not found in the database.")
	}

	// since the schema is set up in the database, this function is done with
	// its task. due to the "defer" statement, the connection will now be closed
	// automatically
	log.Info().Msg("database schema found. closing temporary database connection")
}
