package users

import (
	events "github.com/gokit/cqrskit/examples/users/events"

	"github.com/gokit/cqrskit"
)

var (
	// UserAggregateID represents the unique aggregate id for all events
	// related to the User type. It is the typeName hashed using a md5 sum.
	UserAggregateID = "f7091ac77d9b52a3ec5609891cd9f54f"
)

//*******************************************************************************
// User Event Applier
//*******************************************************************************

// UserEvents implements the necessary logic to apply a
// series of events to a giving User type.
type UserEvents struct {
	Events []cqrskit.Event
}

// Apply embodies the internal logic necessary to apply specific events to a User by
// calling appropriate methods.
func (esv UserEvents) Apply(u *User) error {
	for _, event := range esv.Events {
		switch ev := event.EventData.(type) {
		case UserEmailUpdated:
			return u.HandleUserEmailUpdated(ev)
		case events.UserNameUpdated:
			return u.HandleUserNameUpdated(ev)

		}
	}
	return nil
}
