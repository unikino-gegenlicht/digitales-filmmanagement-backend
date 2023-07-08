package config

import "errors"

// ErrEmptyOpenIdConnectConfig is returned if there are no entries in the configuration
// for the OpenID Connect
var ErrEmptyOpenIdConnectConfig = errors.New("no open id connect endpoints set")

// ErrInvalidDiscoveryURI is returned if the discovery URI of the OpenIDConnect
// configuration is not a valid URI
var ErrInvalidDiscoveryURI = errors.New("invalid discovery uri")

// ErrInvalidDiscoveryResponse is returned if the discovery URI was requested
//
//	successfully, but the response does not match the specification
var ErrInvalidDiscoveryResponse = errors.New("invalid discovery response")

// ErrDiscoveryResponseMissingUserInfo is returned if the discovery request
// was successful, but does not contain the required address for the userinfo
// endpoint
var ErrDiscoveryResponseMissingUserInfo = errors.New("discovery did not disclose userinfo endpoint")

// ErrInvalidUserInfoURI is returned if the user info URI of the OpenIDConnect
// configuration is not a valid URI
var ErrInvalidUserInfoURI = errors.New("invalid userinfo uri")

// ErrInsecureUserInfoURI is returned if the user info URI of the OpenIDConnect
// configuration is valid, but uses http
var ErrInsecureUserInfoURI = errors.New("insecure userinfo uri")

// ErrNoDatabaseHost is returned if the configuration contains no host for the
// MariaDB that is used in this project
var ErrNoDatabaseHost = errors.New("database host not set")

// ErrEmptyDatabaseHost is returned if the configuration contains an empty
// host for the MariaDB that is used in this project
var ErrEmptyDatabaseHost = errors.New("database host is empty")

// ErrNoDatabaseUser is returned if the configuration contains no user for the
// MariaDB that is used in this project
var ErrNoDatabaseUser = errors.New("database user not set")

// ErrEmptyDatabaseUser is returned if the configuration contains an empty
// user for the MariaDB that is used in this project
var ErrEmptyDatabaseUser = errors.New("database user is empty")

// ErrNoDatabaseSpecified is returned if the configuration contains no database name
// for the MariaDB that is used in this project
var ErrNoDatabaseSpecified = errors.New("database name not set")

// ErrEmptyDatabaseSpecified is returned if the configuration contains an empty
// database name for the MariaDB that is used in this project
var ErrEmptyDatabaseSpecified = errors.New("database name is empty")
