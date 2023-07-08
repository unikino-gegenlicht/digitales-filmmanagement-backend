package config

import (
	"fmt"
	"strings"
)

// WpDbConfiguration contains the configuration for the database connection
// to the MariaDB server required for this backend to work. It contains an
// internal boolean to check if it has been validated.
type WpDbConfiguration struct {
	Host      *string `toml:"host"`
	Port      *string `toml:"port"`
	User      *string `toml:"user"`
	Password  *string `toml:"password"`
	Schema    *string `toml:"schema"`
	validated bool
}

// Validate checks if the configuration contains at least the user, host and
// database that should be used by the application. All other variables are
// optional.
func (c *WpDbConfiguration) Validate() error {
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
	if c.Schema == nil {
		return ErrNoDatabaseSpecified
	}
	// now check if the host is empty
	if strings.TrimSpace(*c.Schema) == "" {
		return ErrEmptyDatabaseSpecified
	}

	// now check if the optional port was set. if not, set it to 5432, which is
	// the default port for postgres
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

// BuildDSN returns a connection string for sql.Open. Before
// building the connection string, the configuration needs to be validated with
// Validate. If the configuration is not validated, an empty string will be
// returned
func (c *WpDbConfiguration) BuildDSN() string {
	if c.validated {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			*c.User, *c.Password, *c.Host, *c.Port, *c.Schema)
	}
	return ""

}
