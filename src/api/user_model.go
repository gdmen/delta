package api

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
