package globals

import (
	"database/sql"
	"digitales-filmmanagement-backend/config"
	"github.com/go-chi/httplog"
)

// HttpLogger is the logger used by the code interacting with API requests
var HttpLogger = httplog.NewLogger("management-backend", httplog.Options{JSON: true})

// Configuration is the main configuration for this application. It contains
// every used subsection from the 'configuration.toml' file.
var Configuration config.Configuration

// Database is the shared connection to the MariaDB Database used by the backend
var Database *sql.DB
