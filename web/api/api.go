package api

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/web/api/user"
)

func Handler(c *app.Ctx, nextRoute string) error {
	switch nextRoute {
	case "user":
		c.Next(user.Handler)
		return nil
	default:
		return c.NotFound()
	}
}
