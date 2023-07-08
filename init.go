package main

import (
	"database/sql"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"digitales-filmmanagement-backend/config"
	"digitales-filmmanagement-backend/globals"
	// database driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
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

// this function now loads the prepared errors from the error file and parses
// them into wisdom errors
func init() {
	log.Info().Msg("loading predefined errors")
	file, err := os.Open("./errors.json")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open error configuration file")
	}

	var errors []types.APIError
	err = json.NewDecoder(file).Decode(&errors)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load error configuration file")
	}
	for _, e := range errors {
		e.InferHttpStatusText()
		globals.Errors[e.ErrorCode] = e
	}
	log.Info().Msg("loaded predefined errors")
}

// this function connects the application to the database server and checks if
// the wanted database schema is present on the server. to achieve this, this
// step creates a temporary database connection which connects to a user's
// default database schema. the connection is closed after this step
func init() {
	// get the configuration
	c := globals.Configuration.Database

	// use the configuration to create a dsn without specifying the schema
	dsn := c.BuildDSN()

	// now open a connection to the database
	db, err := sql.Open("postgres", dsn)
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

	// since the schema is set up in the database, this function is done with
	// its task. due to the "defer" statement, the connection will now be closed
	// automatically
	log.Info().Msg("database schema found. closing temporary database connection")
}

// this function now establishes the globally used database connection and
// checks afterward if the required tables are present in the previously
// configured schema
func init() {
	// get the database configuration again
	c := globals.Configuration.Database
	// now build the dsn including the schema name
	dsn := c.BuildDSN()
	// now try to open the connection
	var err error
	globals.Database, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open permanent database connection")
	}

	// now configure the connection pooling used to reduce reconnections
	globals.Database.SetConnMaxLifetime(time.Minute * 3)
	globals.Database.SetMaxOpenConns(50)
	globals.Database.SetMaxIdleConns(50)

	// now ping the database to confirm the connectivity
	err = globals.Database.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to ping database server over permanent connection")
	}

	// now execute the init.sql file to create the needed tables if they are
	// not already present
	log.Info().Msg("checking database for schema and tables. missing objects will be created")
	initQueries, err := dotsql.LoadFromFile("./init.sql")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to load sql file with initial table definitions")
	}
	for name, _ := range initQueries.QueryMap() {
		_, err := initQueries.Query(globals.Database, name)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to create needed database objects")
		}
	}
	log.Info().Msg("connected to postgres")
}

// this function now establishes the globally used database connection and
// checks afterward if the required tables are present in the previously
// configured schema
func init() {
	// get the database configuration again
	wp := globals.Configuration.WordPress
	_ = wp.Validate()
	// now build the dsn including the schema name
	dsn := wp.BuildDSN()
	// now try to open the connection
	var err error
	globals.WpDatabase, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to open permanent database connection")
	}

	// now configure the connection pooling used to reduce reconnections
	globals.WpDatabase.SetConnMaxLifetime(time.Minute * 3)
	globals.WpDatabase.SetMaxOpenConns(50)
	globals.WpDatabase.SetMaxIdleConns(50)

	// now ping the database to confirm the connectivity
	err = globals.WpDatabase.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to ping database server over permanent connection")
	}
	log.Info().Msg("connected to wordpress")
}
