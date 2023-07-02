package routes

import (
	"digitales-filmmanagement-backend/globals"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
	"github.com/blockloop/scan/v2"
	"net/http"
)

func GetAllRegisters(w http.ResponseWriter, r *http.Request) {
	// access the request context and get the global error handler
	ctx := r.Context()
	errorHandler := ctx.Value("nativeErrorChannel").(chan error)
	handledError := ctx.Value("nativeErrorHandled").(chan bool)
	// now try to get all register items from the database
	rows, err := globals.SqlQueries.Query(globals.Database, "get-registers")
	if err != nil {
		// send error to the error handler
		errorHandler <- err
		// wait until the error was handled
		<-handledError
		return
	}
	// now create an array of register items
	var items []types.Register
	// now scan the result rows into the array
	if err = scan.Rows(&items, rows); err != nil {
		// send error to the error handler
		errorHandler <- err
		// wait until the error was handled
		<-handledError
		return
	}
	// now return the available register items
	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		// send error to the error handler
		errorHandler <- err
		// wait until the error was handled
		<-handledError
		return
	}
}
