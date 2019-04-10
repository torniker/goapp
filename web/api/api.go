package api

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/web/api/user"
)

func Handler(c *app.Ctx, nextRoute string) error {
	logger.Infof("nextRoute: %v\n", nextRoute)
	switch nextRoute {
	case "user":
		c.Next(user.Handler)
	default:
		return c.NotFound()
	}
	return nil
}
