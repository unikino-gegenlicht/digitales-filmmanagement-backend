package middleware

import (
	"context"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

// NativeErrorHandler allows the global handling and wrapping of native errors
// occurring in API calls. The function needs the service name as a parameter
// to correctly generate the error code used in the wisdomType.WISdoMError
//
// To access the channel added to the request context in a http handler use
// the following call:
//
//	nativeErrorChannel := r.Context().Value("nativeErrorChannel").(chan error)
//
// To render an error, just send it into the channel using the following syntax:
//
//	nativeErrorChannel<-err
func NativeErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create a new channel
		c := make(chan error)
		a := make(chan bool)
		// now access the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "nativeErrorChannel", c)
		ctx = context.WithValue(ctx, "nativeErrorHandled", a)
		// use a go function to listen to the channel and output the
		// request error to the client using json
		go func() {
			for {
				select {
				case err := <-c:
					e := types.APIError{}
					e.WrapError(err)
					w.Header().Set("Content-Type", "text/json")
					w.WriteHeader(500)
					encodingErr := json.NewEncoder(w).Encode(e)
					if encodingErr != nil {
						log.Error().Err(encodingErr).Msg("unable to send error")
					}
					a <- true
					return
				}
			}
		}()
		// now let the next handler handle the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// APIErrorHandler allows the global handling and wrapping of WISdoM errors
// occurring in API calls. The function needs the mapping of predefined error
// messages using the wisdomType.WISdoMError
//
// To access the channel added to the request context in a http handler use
// the following call:
//
//	wisdomErrorChannel := r.Context().Value("wisdomErrorChannel").(chan string)
//
// To render the error just send the error code into the channel:
//
//	wisdomErrorChannel <- "test"
//
// If an unregistered error code is used a panic will be released in the go
// function handling the actual rendering
func APIErrorHandler(errors map[string]types.APIError) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// create a new channel
			c := make(chan string)
			a := make(chan bool)
			// now access the request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "apiErrorChannel", c)
			ctx = context.WithValue(ctx, "apiErrorHandled", a)
			// use a go function to listen to the channel and output the
			// request error to the client using json
			go func() {
				for {
					select {
					case errorCode := <-c:
						e, errorPresent := errors[errorCode]
						if !errorPresent {
							panic("using unregistered error")
						}
						w.Header().Set("Content-Type", "text/json")
						w.WriteHeader(500)
						encodingErr := json.NewEncoder(w).Encode(e)
						if encodingErr != nil {
							log.Error().Err(encodingErr).Msg("unable to send error")
						}
						a <- true
						return
					}
				}
			}()
			// now let the next handler handle the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
