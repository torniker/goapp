package user

import (
	"github.com/gofrs/uuid"
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/db"
)

func Handler(c *app.Ctx, nextRoute string) error {
	if c.Method() == "POST" {
		// insert user
		return nil
	}
	userID, err := uuid.FromString(nextRoute)
	if err != nil {
		logger.Error(err)
		return c.NotFound()
	}
	udb, err := db.UserByID(c.App.PG(), userID)
	if err != nil {
		logger.Error(err)
		return c.InternalError()
	}
	return c.JSON(udb.Model())
}
