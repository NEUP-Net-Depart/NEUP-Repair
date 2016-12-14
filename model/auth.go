package model

import "github.com/pkg/errors"

type Auth struct {
	ID       int64  `json:"id,omitempty"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	// Keep it simple and safe
}

func PasswordByName(db Storager, name string) (encPass string, err error) {
	sqlStr := "SELECT pass FROM users WHERE name = ?"
	err = db.QueryRowx(sqlStr, name).Scan(&encPass)
	if err != nil {
		err = errors.Wrap(err, "get pass error, query error")
		return
	}
	return
}

func (au Auth) Add(db Storager) (err error) {
	sqlStr := "INSERT INTO users (`name`, `pass`)VALUES(?, ?)"
	_, err = db.Exec(sqlStr, au.UserName, au.Password)
	if err != nil {
		err = errors.Wrap(err, "add auth error, SQL insert error")
		return
	}
	return
}
