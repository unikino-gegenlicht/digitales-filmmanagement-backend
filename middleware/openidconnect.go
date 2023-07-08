package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"digitales-filmmanagement-backend/config"
	"digitales-filmmanagement-backend/globals"
)

func UserInfo(c config.OpenIdConnectConfiguration) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// access the current request context
			ctx := request.Context()
			// now get the global error handler
			nativeErrorChannel := ctx.Value("nativeErrorChannel").(chan error)
			handledNativeError := ctx.Value("nativeErrorHandled").(chan bool)
			apiErrorHandler := ctx.Value("apiErrorChannel").(chan string)
			handledApiError := ctx.Value("apiErrorHandled").(chan bool)
			// now access the request's headers and get the authorization header value
			authHeaderValue := strings.TrimSpace(request.Header.Get("Authorization"))
			// now check if the header actually does contain anything
			if authHeaderValue == "" {
				apiErrorHandler <- "MISSING_AUTHORIZATION_HEADER"
				<-handledApiError
				return
			}
			// since some value was set in the authorization header, now build the request for the userinfo endpoint of
			// the OpenIDConnect server
			userinfoRequest, err := http.NewRequest("GET", *c.UserInfoEndpoint, nil)
			if err != nil {
				nativeErrorChannel <- err
				<-handledNativeError
				return
			}
			userinfoRequest.Header.Set("Authorization", authHeaderValue)

			// now execute the request
			userinfoResponse, err := globals.HttpClient.Do(userinfoRequest)
			if err != nil {
				nativeErrorChannel <- err
				<-handledNativeError
				return
			}
			// now check the response code
			switch userinfoResponse.StatusCode {
			case 200:
				// now parse the user info into a map
				userInfo := make(map[string]interface{})
				err = json.NewDecoder(userinfoResponse.Body).Decode(&userInfo)
				if err != nil {
					nativeErrorChannel <- err
					<-handledNativeError
					return
				}

				// now set the full name of the user and the groups into the
				// context
				ctx = context.WithValue(ctx, "user", userInfo["given_name"])
				break
			case 401:
				apiErrorHandler <- "UNAUTHORIZED"
				<-handledApiError
				return
			case 403:
				apiErrorHandler <- "FORBIDDEN"
				<-handledApiError
				return
			default:
				nativeErrorChannel <- errors.New("unexpected response code during authentication validation")
				<-handledNativeError
				return
			}

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
