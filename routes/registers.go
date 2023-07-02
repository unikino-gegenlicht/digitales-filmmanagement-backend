package routes

import (
	"digitales-filmmanagement-backend/globals"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func NewRegisterTransaction(w http.ResponseWriter, r *http.Request) {
	// access the request context and get the global error handler
	ctx := r.Context()
	nativeErrorHandler := ctx.Value("nativeErrorChannel").(chan error)
	handledNativeError := ctx.Value("nativeErrorHandled").(chan bool)
	apiErrorHandler := ctx.Value("apiErrorChannel").(chan string)
	handledApiError := ctx.Value("apiErrorHandled").(chan bool)

	// now get the full name of the person responsible for the transaction
	responsiblePerson := ctx.Value("user").(string)

	// now first get the register id from the request
	registerId := chi.URLParam(r, "registerId")
	// now check if the uuid is valid
	if _, err := uuid.Parse(registerId); err != nil {
		apiErrorHandler <- "INVALID_REGISTER_UUID"
		<-handledApiError
		return
	}

	// now try and parse the request body
	var registerTransaction types.RegisterTransaction
	if err := json.NewDecoder(r.Body).Decode(&registerTransaction); err != nil {
		switch err := err.(type) {
		case *json.SyntaxError:
			log.Warn().Err(err).Str("error", "INVALID_JSON").Msg("received invalid json payload")
			apiErrorHandler <- "INVALID_JSON"
			<-handledApiError
			return
		case *json.UnmarshalTypeError:
			log.Warn().Err(err).Str("error", "INVALID_TRANSACTION").Msg("received invalid json payload")
			apiErrorHandler <- "INVALID_TRANSACTION"
			<-handledApiError
			return
		case *json.InvalidUnmarshalError:
			log.Error().Err(err).Msg("invalid unmarshal argument")
			nativeErrorHandler <- err
			<-handledNativeError
			return
		default:
			log.Error().Err(err).Msg("unable to unmarshal into struct")
			nativeErrorHandler <- err
			<-handledNativeError
			return
		}
	}

	// now build a transaction that can be inserted into the database
	transaction := types.Transaction{
		Title:       registerTransaction.Title,
		Description: &registerTransaction.Description,
		Amount:      registerTransaction.Total,
		By:          responsiblePerson,
		Register:    registerId,
	}
	_, err := globals.SqlQueries.Exec(globals.Database, "insert-transaction",
		transaction.Title, transaction.Description, transaction.Amount, transaction.By, transaction.Register)
	if err != nil {
		log.Error().Err(err).Msg("error while inserting transaction")
		nativeErrorHandler <- err
		<-handledNativeError
		return
	}
	// now insert the statistics
	for articleName, articleCount := range registerTransaction.Articles {
		_, err = globals.SqlQueries.Exec(globals.Database, "insert-article-sale",
			articleName, articleCount)
		if err != nil {
			log.Error().Err(err).Msg("error while inserting article statistics")
			nativeErrorHandler <- err
			<-handledNativeError
			return
		}
	}

	// now report back that the transaction was completely stored in the database
	w.WriteHeader(http.StatusCreated)
	return
}
