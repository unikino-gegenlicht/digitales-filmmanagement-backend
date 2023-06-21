package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"digitales-filmmanagement-backend/config"
	"digitales-filmmanagement-backend/globals"
	httpTypes "digitales-filmmanagement-backend/types/http"
)

func UserInfo(c config.OpenIdConnectConfiguration) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// access the current request context
			ctx := request.Context()
			// now access the request's headers and get the authorization header value
			authHeaderValue := strings.TrimSpace(request.Header.Get("Authorization"))
			// now check if the header actually does contain anything
			if authHeaderValue == "" {
				e := httpTypes.ErrorMessage{
					Error:   "NO_AUTH_HEADER_SET",
					Message: "THe request lacks the `Authorization` header",
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
				e := httpTypes.ErrorMessage{
					Error:   "INTERNAL_ERROR",
					Message: "Error while building request for userinfo: " + err.Error(),
				}
				writer.WriteHeader(401)
				_ = json.NewEncoder(writer).Encode(e)
				return
			}
			userinfoRequest.Header.Set("Authorization", authHeaderValue)

			// now execute the request
			userinfoResponse, err := globals.HttpClient.Do(userinfoRequest)
			if err != nil {
				e := httpTypes.ErrorMessage{
					Error:   "INTERNAL_ERROR",
					Message: "Error while requesting userinfo: " + err.Error(),
				}
				writer.Header().Set("Content-Type", "text/json")
				writer.WriteHeader(401)
				_ = json.NewEncoder(writer).Encode(e)
				return
			}
			// now parse the user info into a map
			userInfo := make(map[string]interface{})
			json.NewDecoder(userinfoResponse.Body).Decode(&userInfo)

			// now set the full name of the user and the groups into the
			// context
			ctx = context.WithValue(ctx, "username", userInfo["preferred_username"])

			fmt.Println(userInfo)
		})
	}
}
