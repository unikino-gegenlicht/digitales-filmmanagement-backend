package globals

import (
	"github.com/go-chi/httplog"
)

// HttpLogger is the logger used by the code interacting with API requests
var HttpLogger = httplog.NewLogger("management-backend", httplog.Options{JSON: true})
