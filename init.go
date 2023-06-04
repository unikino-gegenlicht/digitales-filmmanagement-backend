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

// This function connects the application to the MariaDB server and checks if
// all needed tables are present in the database. If not, the tables will be
// created by this function
func init() {
	// get the mariadb configuration
	dbConfig := globals.Configuration.Database
	// now get the connection string
	dsn := dbConfig.BuildConnectionString()
	log.Debug().Str("dsn", dsn).Msg("built data source name")
	// now open the connection to the database
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	// now set some connection properties for better usability/less errors
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(50)
	conn.SetMaxIdleConns(50)
	// now ping the database again, to verify the connection
	err = conn.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("ping failed but initial connection successful")
	}
	// since all connections have been successful, set the global connection
	// to the one just created
	globals.Database = conn
	// TODO: implement database checks
	//	 the checks shall include:
	//	    - existence of database schema
	//	    - existence of tables
	//		- definition of tables
}
