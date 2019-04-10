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
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Model converts schema User to model User
func (udb *User) Model() model.User {
	user := model.User{
		ID:        udb.ID,
		Username:  udb.Username,
		CreatedAt: udb.CreatedAt,
	}
	return user
}
