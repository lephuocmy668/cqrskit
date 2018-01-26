package usercqrs

import (
	"github.com/gokit/cqrskit/examples/users"

	"github.com/gokit/cqrskit/internal/cqrs"
)

var (
	// UserAggregateID represents the unique aggregate id for all events
	// related to the User type. It is the typeName hashed using a md5 sum.
	UserAggregateID = "f7091ac77d9b52a3ec5609891cd9f54f"
)

//*******************************************************************************
// Event Handler
//*******************************************************************************

// UserEvents implements the necessary logic to apply a
// series of events to a giving User type.
type UserEvents struct {
	Events []cqrs.Event
}

// Apply embodies the internal logic necessary to apply specific events to a User by
// calling appropriate methods.
func (esv UserEvents) Apply(U *users.User) error {

	return nil
}
