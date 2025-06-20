package models

import (
	"log"
	"time"

	"github.com/SabianF/ghasp/src/common/data/sources"
	"github.com/SabianF/ghasp/src/common/domain/entities"
	"github.com/gofrs/uuid"
)

type UserModel interface {
	Create() (entities.User, error)
	Query() (entities.User, error)
	Update() (entities.User, error)
	User() entities.User
	UserModel() userModel
}

type userModel struct {
	Id string
	Created_at string
	Updated_at string
	Name_first string
	Name_last string
	Email string
}

func CreateUserModel(user entities.User) UserModel {
	return &userModel{
		Id: uuid.Must(uuid.NewV7()).String(),
		Name_first: user.User().Name_first,
		Name_last: user.User().Name_last,
		Email: user.User().Email,
	}
}

func (model *userModel) User() entities.User {
	user, err := entities.CreateUserFromModel(
		model.Id,
		model.Created_at,
		model.Updated_at,
		model.Name_first,
		model.Name_last,
		model.Email,
	)
	if (err != nil) {
		log.Println("failed to get User from UserModel: " + err.Error())
		return nil
	}

	return user
}

func (model *userModel) UserModel() userModel {
	return *model
}

func (model *userModel) Create() (entities.User, error) {
	modelWithTimestamps := userModel{
		Id: model.Id,
		Created_at: time.Now().Format(time.RFC3339Nano),
		Updated_at: time.Now().Format(time.RFC3339Nano),
		Name_first: model.Name_first,
		Name_last: model.Name_last,
		Email: model.Email,
	}

	createdUser, err := sources.Db.CreateUser(modelWithTimestamps.User())

	return createdUser, err
}

func (model *userModel) Query() (entities.User, error) {
	return sources.Db.GetUser(model.User().User().Id)
}

func (model *userModel) Update() (entities.User, error) {
	return sources.Db.UpdateUser(model.User())
}

func GetUserFieldNames() []string {
	return []string{
		"First name",
		"Last name",
		"Email",
	}
}

func GetAllUsers() ([]entities.User, error) {
	return sources.Db.GetAllUsers()
}
