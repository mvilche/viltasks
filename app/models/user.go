package models

import (
	"errors"

	"github.com/revel/config"
	"github.com/revel/revel"
)

type User struct {
	Username string
	Password string
}

func GetUser() (User, error) {

	var u User
	c, err := config.ReadDefault("./conf/app.conf")
	if err != nil {
		revel.AppLog.Error("", err)
		return u, err
	}

	var user, _ = c.String("auth", "auth.username")
	var pass, _ = c.String("auth", "auth.password")

	if len(user) > 0 && len(pass) > 0 {
		u.Username = user
		u.Password = pass
		return u, nil

	} else {

		err := errors.New("Password or username empty")
		return u, err
	}
}
