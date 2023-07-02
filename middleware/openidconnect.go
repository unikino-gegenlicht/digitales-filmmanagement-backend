package middleware

import (
	"context"
	"digitales-filmmanagement-backend/types"
	"encoding/json"
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
			// now access the request's headers and get the authorization header value
			authHeaderValue := strings.TrimSpace(request.Header.Get("Authorization"))
			// now check if the header actually does contain anything
			if authHeaderValue == "" {
				e := types.APIError{
					ErrorCode:        "NO_AUTH_HEADER_SET",
					ErrorTitle:       "No Authorization Header set",
					ErrorDescription: "The request needs to have the Authorization header set to allow access to the API",
					HttpStatusCode:   401,
					HttpStatusText:   http.StatusText(401),
				}
				writer.Header().Set("Content-Type", "text/json")
				writer.WriteHeader(401)
				json.NewEncoder(writer).Encode(e)
				return
			}
			// since some value was set in the authorization header, now build the request for the userinfo endpoint of
			// the OpenIDConnect server
			userinfoRequest, err := http.NewRequest("GET", *c.UserInfoEndpoint, nil)
			if err != nil {
				nativeErrorChannel <- err
				return
			}
			userinfoRequest.Header.Set("Authorization", authHeaderValue)

			// now execute the request
			userinfoResponse, err := globals.HttpClient.Do(userinfoRequest)
			if err != nil {
				nativeErrorChannel <- err
				return
			}
			// now parse the user info into a map
			userInfo := make(map[string]interface{})
			err = json.NewDecoder(userinfoResponse.Body).Decode(&userInfo)
			if err != nil {
				nativeErrorChannel <- err
				return
			}

			// now set the full name of the user and the groups into the
			// context
			ctx = context.WithValue(ctx, "user", userInfo["given_name"])

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
