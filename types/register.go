package types

// Register reflects a cash register stored in the database
type Register struct {
	// ID contains the UUID used to identify the register in API calls
	ID *string `json:"id" db:"id"`
	// Name contains the name used to identify the register in frontend applications
	Name string `json:"name"`
	// Description contains a optional
	Description *string `json:"description"`
}
