package sources

import (
	"fmt"
	"log"

	"github.com/SabianF/ghasp/src/common/domain/entities"
)

type Database struct {
	users []entities.User
}

var db = Database{}

func InitDb() {
	log.Println("Opening DB connection...")

	db.users = []entities.User{}

	for i := 0; i < 10; i++ {
		newUser, err := entities.NewUser(
			"John",
			fmt.Sprintf("Smith %v", i + 1),
			fmt.Sprintf("john.smith.%v@email.com", i + 1),
		)

		if (err != nil) {
			log.Println(err)
			continue
		}

		db.users = append(db.users, newUser)
	}

	log.Println("Successfully opened DB connection.")

}

func CloseDb() {
	log.Println("Closing database connection...")

	db.users = nil

	log.Println("Done closing database connection.")
}
