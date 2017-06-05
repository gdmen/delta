package api

import (
	"errors"
)

const (
	UsernameUnavailableUserErrMsg = "That username isn't available"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

// For this demo, we're storing the user list in memory
// We also have some users predefined.
// In a real application, this list will most likely be fetched
// from a database. Moreover, in production settings, you should
// store passwords securely by salting and hashing them instead
// of using them as we're doing in this demo
var userList = []User{
	User{Username: "username", Password: "pass"},
	User{Username: "user2", Password: "pass2"},
	User{Username: "user3", Password: "pass3"},
}

// Register a new user with the given username and password
func registerNewUser(username, password string) (*User, error) {
	if !isUsernameAvailable(username) {
		return nil, errors.New(UsernameUnavailableUserErrMsg)
	}

	u := User{Username: username, Password: password}

	userList = append(userList, u)

	return &u, nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	for _, u := range userList {
		if u.Username == username {
			return false
		}
	}
	return true
}
