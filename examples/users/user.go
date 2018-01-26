package users

import "github.com/gokit/cqrskit/examples/users/events"

//@escqrs
type User struct {
	Version  int
	Email    string
	Username string
}

type UserEmailUpdated struct {
	New string
}

func (u *User) HandleUserEmailUpdated(ev UserEmailUpdated) error {
	return nil
}

func (u *User) HandleUserNameUpdated(ev events.UserNameUpdated) error {
	return nil
}

//@escqrs-method-skip
func (u *User) HandleUserRackUpdated(ev UserEmailUpdated) error {
	return nil
}
