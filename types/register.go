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

type RegisterTransaction struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Total       float64        `json:"amount"`
	Articles    map[string]int `json:"articleCounts"`
}
