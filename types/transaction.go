package types

type Transaction struct {
	// Title contains the title of the transaction
	Title string `json:"title" db:"title"`
	// Description contains a more in depth description of the transaction
	Description *string `json:"description" db:"description"`
	// Amount contains the amount of the transaction in euros
	Amount float64 `json:"amount" db:"amount"`
	// By contains the full name of the person responsible for this transaction
	By string `json:"by" db:"by"`
	// Register contains the register in which th transaction took place
	Register string `json:"register" db:"register"`
	// storedInDb contains a boolean indicator to stop writing the transaction
	// multiple times into the database
	storedInDb bool
}
