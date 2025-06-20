package sources

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/SabianF/ghasp/src/common/domain/entities"
	"github.com/gofrs/uuid"
)

type Database struct {
	users []entities.User
}

var Db = Database{}

func InitDb() {
	log.Println("Opening DB connection...")

	Db.users = []entities.User{}

	for i := 0; i < 10; i++ {
		newUser, err := entities.CreateUserFromModel(
			uuid.Must(uuid.NewV7()).String(),
			time.Now().Format(time.RFC3339Nano),
			time.Now().Format(time.RFC3339Nano),
			"John",
			fmt.Sprintf("Smith %v", i + 1),
			fmt.Sprintf("john.smith.%v@email.com", i + 1),
		)

		if (err != nil) {
			log.Println(err)
			continue
		}

		Db.users = append(Db.users, newUser)
	}

	log.Println("Successfully opened DB connection.")

}

func CloseDb() {
	log.Println("Closing database connection...")

	Db.users = nil

	log.Println("Done closing database connection.")
}

func (db Database) GetAllUsers() ([]entities.User, error) {
	if (Db.users == nil) {
		return nil, errors.New("users database is nil")
	}

	return Db.users, nil
}

func (db Database) CreateUser(user entities.User) (entities.User, error) {
	if (user.User().Id != "") {
		_, _, err := findUser(user.User().Id)
		if (err == nil) {
			return nil, errors.New("user already exists")
		}
		if (err.Error() != "user not found") {
			return nil, err
		}
	}

	Db.users = append(Db.users, user)

	createdUser, _, err := findUser(user.User().Id)
	if (err != nil) {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	return createdUser, nil
}

func (db Database) GetUser(id string) (entities.User, error) {
	foundUser, _, err := findUser(id)
	if (err != nil) {
		return nil, err
	}

	return foundUser, nil
}

func (db Database) UpdateUser(user entities.User) (entities.User, error) {
	foundUserBeforeChanges, foundUserIndex, err := findUser(user.User().Id)
	if (err != nil) {
		return nil, errors.New(err.Error())
	}

	log.Println(Db.users[foundUserIndex].User())
	Db.users[foundUserIndex] = user
	log.Println(Db.users[foundUserIndex].User())

	foundUserAfterChanges, _, err := findUser(user.User().Id)
	if (err != nil) {
		return nil, errors.New("failed to update user: " + err.Error())
	}
	if (foundUserAfterChanges == foundUserBeforeChanges) {
		return nil, errors.New("failed to update user: data in database is still the same")
	}

	return foundUserAfterChanges, nil
}

func findUser(id string) (entities.User, int, error) {
	if (id == "") {
		return nil, -1, errors.New("valid id not provided")
	}

	for i, user := range Db.users {
		if (user.User().Id == id) {
			return user, i, nil
		}
	}

	return nil, -1, errors.New("user not found")
}
