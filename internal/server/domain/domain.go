package domain

import "fmt"

type User struct {
	ID       string
	Username string
	Password string
	Email    string
}

func (u User) String() string {
	return fmt.Sprintf("User ID: %s, Username: %s, Email: %s", u.ID, u.Username, u.Email)
}
