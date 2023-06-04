package dbTypes

import "errors"

// ErrNoDbAction indicates that no action was supplied to the function requiring
// it.
var ErrNoDbAction = errors.New("no db action supplied")

// ErrInvalidDbAction indicates that the action supplied to a function is not valid
// for this action
var ErrInvalidDbAction = errors.New("invalid db action supplied")
