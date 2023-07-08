package types

// RegisterItem represents an item which can be sold
type RegisterItem struct {
	// ID contains the UUID of the item
	ID *string `json:"id" db:"id"`
	// Name contains the name for the item
	Name string `json:"name" db:"name"`
	// Price contains the price of the item
	Price float64 `json:"price" db:"price"`
	// Icon contains a string pointing to an icon which is displayed next to an
	// item
	Icon string `json:"icon" db:"icon"`
}
