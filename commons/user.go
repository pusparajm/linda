package commons

type User struct {
	Username string
	Nickname string
}

func NewUser(username, nickname string) *User {
	u := new(User)
	u.Username = username
	u.Nickname = nickname
	return u
}
