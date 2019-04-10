package user

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/db"
	"github.com/torniker/goapp/schema"
)

func Handler(c *app.Ctx, nextRoute string) error {
	if c.Method() == "POST" {
		return handleInsert(c, nextRoute)
	}
	userID, err := uuid.FromString(nextRoute)
	if err != nil {
		logger.Warn(err)
		return c.NotFound()
	}
	user, err := db.UserByID(c.App.PG(), userID)
	if err != nil {
		logger.Error(err)
		return c.InternalError()
	}
	if user == nil {
		return c.NotFound()
	}
	return c.JSON(user.Model())
}

type userInsertRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleInsert(c *app.Ctx, nextRoute string) error {
	decoder := json.NewDecoder(c.RequestBody())
	var uir userInsertRequest
	err := decoder.Decode(&uir)
	if err != nil {
		logger.Error(err)
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		logger.Error(err)
		return err
	}
	userDB := schema.User{
		ID:        id,
		Username:  uir.Username,
		Password:  uir.Password,
		CreatedAt: time.Now(),
	}
	err = db.UserInsert(c.App.PG(), userDB)
	if err != nil {
		return err
	}
	return c.JSON(userDB.Model())
}
