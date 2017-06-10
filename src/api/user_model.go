package api

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	UsernameUnavailableErrMsg  = "That username isn't available"
	CouldNotHashPasswordErrMsg = "Couldn't hash password"
	CouldNotInsertUserErrMsg   = "Couldn't add user to database"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"-" form:"password"`
}

const (
	CreateUserTableSQL = `
create table users (
	id int auto_increment primary key,
	username varchar not null,
	password varchar not null,
	constraint username unique (username)
);`
	InsertUserSQL = `
insert into users(username, password) values(?, ?);`
)

func registerNewUser(db *sql.DB, user *User) (*User, error) {
	// Salt and hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %v", CouldNotHashPasswordErrMsg, err))
	}

	// Write to DB
	result, err := db.Exec(InsertUserSQL, user.Username, hashedPass)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New(UsernameUnavailableErrMsg)
		}
		return nil, errors.New(fmt.Sprintf("%s: %v", CouldNotInsertUserErrMsg, err))
	}
	id, err := result.LastInsertId()
	if err != nil {
		// If the db doesn't support LastInsertId(), throw an error for now
		panic(fmt.Sprintf("Unsupported database caused: %v", err))
	}
	user.Id = id
	user.Password = ""
	return user, nil
}
