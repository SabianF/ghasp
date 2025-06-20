package entities

import (
	"errors"
	"net/mail"
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

	if (name_first == "") {
		err := errors.New("missing first name")
		return nil, err
	}

	if (name_last == "") {
		err := errors.New("missing last name")
		return nil, err
	}

	if (email == "") {
		err := errors.New("missing email")
		return nil, err
	}

	parsedEmail, err := mail.ParseAddress(email)
	if (err != nil) {
		err = errors.Join(errors.New("invalid email: "), err)
		return nil, err
	}

	newUser := user{
		Name_first: name_first,
		Name_last: name_last,
		Email: parsedEmail.Address,
	}

	return &newUser, nil
}

func CreateUserFromModel(
	id string,
	created_at string,
	updated_at string,
	name_first string,
	name_last string,
	email string,
) (User, error) {

	created_at_parsed, errCreatedAt := time.Parse(time.RFC3339Nano, created_at)
	if (errCreatedAt != nil) {
		return nil, errCreatedAt
	}

	updated_at_parsed, errUpdatedAt := time.Parse(time.RFC3339Nano, updated_at)
	if (errUpdatedAt != nil) {
		return nil, errUpdatedAt
	}

	newUser := user{
		Id: id,
		Created_at: created_at_parsed,
		Updated_at: updated_at_parsed,
		Name_first: name_first,
		Name_last: name_last,
		Email: email,
	}

	return &newUser, nil
}

func (user *user) User() user {
	return *user
}
