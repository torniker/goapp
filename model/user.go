package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User is struct for moving user data between micro-services and frontend
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UpdatedAt time.Time `json:"updated_at"`
}
