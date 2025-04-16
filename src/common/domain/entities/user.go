package models

import (
	"errors"
	"time"
)

type User interface {
	User() user
}

type user struct {
	Id string
	Created_at time.Time
	Updated_at time.Time
	Name_first string
	Name_last string
	Email string
}

func NewUser(
	name_first string,
	name_last string,
	email string,
) (User, error) {

	if (name_first == "" || name_last == "" || email == "") {
		err := errors.New("new User is missing name or email")
		return nil, err
	}

	userResult := user{
		Name_first: name_first,
		Name_last: name_last,
		Email: email,
	}

	return &userResult, nil
}

func (user *user) User() user {
	return *user
}
