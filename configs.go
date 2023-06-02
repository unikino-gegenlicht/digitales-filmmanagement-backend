package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
func (c OpenIdConnectConfiguration) Validate() error {
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

// MariaDbConfiguration contains the configuration for the database connection
// to the MariaDB server required for this backend to work. It contains an
// internal boolean to check if it has been validated.
type MariaDbConfiguration struct {
	Host      *string `toml:"host"`
	Port      *string `toml:"port"`
	User      *string `toml:"user"`
	Password  *string `toml:"password"`
	Database  *string `toml:"database"`
	validated bool
}

// Validate checks if the configuration contains at least the user, host and
// database that should be used by the application. All other variables are
// optional.
func (c MariaDbConfiguration) Validate() error {
	// check if the host was set in the configuration
	if c.Host == nil {
		return ErrNoDatabaseHost
	}
	// now check if the host is empty
	if strings.TrimSpace(*c.Host) == "" {
		return ErrEmptyDatabaseHost
	}
	// check if the host was set in the configuration
	if c.User == nil {
		return ErrNoDatabaseUser
	}
	// now check if the host is empty
	if strings.TrimSpace(*c.User) == "" {
		return ErrEmptyDatabaseUser
	}
	// check if the host was set in the configuration
	if c.Database == nil {
		return ErrNoDatabaseSpecified
	}
	// now check if the host is empty
	if strings.TrimSpace(*c.Database) == "" {
		return ErrEmptyDatabaseSpecified
	}

	// now check if the optional port was set. if not, set it to 3306, which is
	// the default port for MariaDb
	if c.Port == nil {
		defaultPort := "3306"
		c.Port = &defaultPort
	}
	// the same procedure is done if the port is empty
	if strings.TrimSpace(*c.Port) == "" {
		defaultPort := "3306"
		c.Port = &defaultPort
	}

	// now check if the optional password was set. if not, set it to an empty
	// string to allow the successful creation of the connection string
	if c.Password == nil {
		defaultPassword := ""
		c.Password = &defaultPassword
	}
	// the same procedure is done if the password is empty
	if strings.TrimSpace(*c.Password) == "" {
		defaultPassword := ""
		c.Password = &defaultPassword
	}

	// since no errors occurred, return nil to indicate that no error occurred
	// and set the validation indicator to true
	c.validated = true
	return nil

}

// BuildConnectionString returns a connection string for sql.Open. Before
// building the connection string, the configuration needs to be validated with
// Validate. If the configuration is not validated, an empty string will be
// returned
func (c MariaDbConfiguration) BuildConnectionString() string {
	if c.validated {
		return fmt.Sprintf("%s:%s@%s:%s/%s?parseTime=true",
			*c.User, *c.Password, *c.Host, *c.Port, *c.Database)
	}
	return ""

}

// Configuration contains all sub-configurations and puts them into one struct
// to allow the parsing of a configuration.toml file which needs to be supplied
// to the documented location (see README or INSTALLATION)
type Configuration struct {
	OIDC     OpenIdConnectConfiguration `toml:"oidc"`
	Database MariaDbConfiguration       `toml:"database"`
}
