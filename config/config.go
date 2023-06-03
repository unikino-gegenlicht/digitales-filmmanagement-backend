package config

// Configuration contains all sub-configurations and puts them into one struct
// to allow the parsing of a configuration.toml file which needs to be supplied
// to the documented location (see README or INSTALLATION)
type Configuration struct {
	OIDC     OpenIdConnectConfiguration `toml:"oidc"`
	Database MariaDbConfiguration       `toml:"database"`
}
