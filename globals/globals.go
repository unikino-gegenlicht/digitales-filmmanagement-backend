package globals

import (
	"database/sql"

	"github.com/go-chi/httplog"
	"github.com/qustavo/dotsql"

	"digitales-filmmanagement-backend/config"
)

// HttpLogger is the logger used by the code interacting with API requests
var HttpLogger = httplog.NewLogger("management-backend", httplog.Options{JSON: true})

// Configuration is the main configuration for this application. It contains
// every used subsection from the 'configuration.toml' file.
var Configuration config.Configuration

// Database is the shared connection to the MariaDB Database used by the backend
var Database *sql.DB

// SqlQueries contains the loaded sql queries from `queries.sql`
var SqlQueries *dotsql.DotSql
