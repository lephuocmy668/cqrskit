package users

type UpdateUserEmail struct {
	New string
}

type EmailUpdatedEvent struct {
	Old string
	New string
}

//@escqrs
//@CommandEvent(Command => UpdateUserEmail, Event => EmailUpdatedEvent)
type User struct {
	Version  int
	Email    string
	Username string
}

func (u *User) ApplyEmailUpdatedEvent(ev UpdateUserEmail) error {
	u.Email = ev.New
	return nil
}

func (u User) HandleUpdateUserEmail(cmd UpdateUserEmail) EmailUpdatedEvent {
	return EmailUpdatedEvent{
		Old: u.Email,
		New: cmd.New,
	}
}
