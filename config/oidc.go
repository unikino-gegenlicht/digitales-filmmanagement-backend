package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// OpenIdConnectConfiguration contains the configuration for the authentication
// of requests for this backend. Either the DiscoveryEndpoint or the
// UserInfoEndpoint need to be set for a valid configuration. If both values
// are set, the UserInfoEndpoint takes precedence.
//
// When using the DiscoveryEndpoint, the function Validate takes the supplied
// URI checks it and sets the UserInfoEndpoint. More information is documented
// in the function itself.
type OpenIdConnectConfiguration struct {

	// DiscoveryEndpoint contains a URI pointing to an Open ID Connect Discovery
	// 1.0 compliant endpoint of an authorization server
	//
	// See also: https://openid.net/specs/openid-connect-discovery-1_0.html
	DiscoveryEndpoint *string `toml:"discoveryEndpoint"`

	// UserInfoEndpoint contains a URI pointing to an Open ID Connect 1.0
	// compliant userinfo endpoint of an authorization server.
	//
	// See also: https://openid.net/specs/openid-connect-core-1_0.html#UserInfo
	UserInfoEndpoint *string `toml:"userInfoEndpoint"`
}

// Validate checks if either the user info endpoint was set in the open id
// connect configuration. If not, it checks if the discovery endpoint was
// set and requests its information and builds the user information endpoint
// address from the response.
//
// If something fails or otherwise does not work, an error will be returned
func (c *OpenIdConnectConfiguration) Validate() error {
	// check if at least one option was set
	if c.DiscoveryEndpoint == nil && c.UserInfoEndpoint == nil {
		return ErrEmptyOpenIdConnectConfig
	}
	// now check if the user info endpoint string contains a valid url
	if c.UserInfoEndpoint == nil {
		// since the user info endpoint was not set, use the discovery endpoint
		// to get the userinfo endpoint.

		// check if the provided discovery endpoint is a valid uri
		uri, err := url.Parse(*c.DiscoveryEndpoint)
		if err != nil {
			return errors.Join(ErrInvalidDiscoveryURI, err)
		}

		// since the uri is valid, check if the uri either uses http or https
		// for the request
		if uri.Scheme != "http" && uri.Scheme != "https" {
			return ErrInvalidDiscoveryURI
		}

		// now make the request
		rawDiscoveryResponse, err := http.Get(*c.DiscoveryEndpoint)
		if err != nil {
			return err
		}

		// now parse the request into a map
		discoveryResponse := make(map[string]interface{})
		err = json.NewDecoder(rawDiscoveryResponse.Body).Decode(&discoveryResponse)
		if err != nil {
			return errors.Join(ErrInvalidDiscoveryResponse, err)
		}

		// now check if the response contains a value for the userinfo endpoint
		userInfoEndpoint, isSet := discoveryResponse["userinfo_endpoint"]
		if !isSet {
			return ErrDiscoveryResponseMissingUserInfo
		}
		userInfoEndpointUri, isString := userInfoEndpoint.(string)
		if !isString {
			return ErrInvalidDiscoveryResponse
		}
		c.UserInfoEndpoint = &userInfoEndpointUri
	}

	endpointUrl, err := url.Parse(*c.UserInfoEndpoint)
	if err != nil {
		return errors.Join(ErrInvalidUserInfoURI, err)
	}
	// now check that the scheme is https since the specification requires the
	// usage of https
	if endpointUrl.Scheme != "https" {
		return ErrInsecureUserInfoURI
	}

	return nil
}
