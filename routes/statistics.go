package routes

import (
	"digitales-filmmanagement-backend/globals"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
	"errors"
	"github.com/blockloop/scan/v2"
	"github.com/go-chi/chi/v5"
	"time"

	"github.com/ggicci/httpin"
	"net/http"
)

func StatisticsRouter() http.Handler {
	r := chi.NewRouter()
	r.With(httpin.NewInput(types.StatisticsRequestInput{})).
		Get("/items", itemStatistics)
	return r
}

func itemStatistics(w http.ResponseWriter, r *http.Request) {
	// access the request context and get the global error handler
	ctx := r.Context()
	nativeErrorHandler := ctx.Value("nativeErrorChannel").(chan error)
	handledNativeError := ctx.Value("nativeErrorHandled").(chan bool)
	// now parse the parameters
	parameters := ctx.Value(httpin.Input).(*types.StatisticsRequestInput)
	var from, until time.Time
	switch {
	case parameters.From == nil && parameters.Until == nil:
		from = time.Now().Add(-24 * time.Hour)
		until = time.Now()
		break
	case parameters.From != nil && parameters.Until == nil:
		from = time.Unix(*parameters.From, 0)
		until = time.Now()
		break
	case parameters.From == nil && parameters.Until != nil:
		from = time.Time{}
		until = time.Unix(*parameters.Until, 0)
		break
	case parameters.From != nil && parameters.Until != nil:
		from = time.Unix(*parameters.From, 0)
		until = time.Unix(*parameters.Until, 0)
		break
	default:
		nativeErrorHandler <- errors.New("unmatched case")
		<-handledNativeError
		return
	}

	rows, err := globals.SqlQueries.Query(globals.Database, "get-article-statistics", from, until)
	if err != nil {
		nativeErrorHandler <- err
		<-handledNativeError
		return
	}

	var statistics []types.ArticleStatistic
	// now parse the rows
	err = scan.Rows(&statistics, rows)
	if err != nil {
		nativeErrorHandler <- err
		<-handledNativeError
		return
	}

	// now return the statistics
	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(statistics)
	if err != nil {
		// send error to the error handler
		nativeErrorHandler <- err
		// wait until the error was handled
		<-handledNativeError
		return
	}
}
