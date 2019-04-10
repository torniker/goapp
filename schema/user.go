package schema

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/torniker/goapp/model"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Model converts schema User to model User
func (udb *User) Model() model.User {
	user := model.User{
		ID:        udb.ID,
		Username:  udb.Username,
		FirstName: udb.FirstName,
		LastName:  udb.LastName,
		CreatedAt: udb.CreatedAt,
	}
	return user
}
