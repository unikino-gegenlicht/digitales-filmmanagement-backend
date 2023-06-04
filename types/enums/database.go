package dbEnums

// DbAction reflects actions that can be performed on database objects
type DbAction int

const (
	DB_UPDATE DbAction = iota
	DB_INSERT
)
