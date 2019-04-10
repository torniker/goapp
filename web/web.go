package web

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/web/api"
)

func Handler(c *app.Ctx, nextRoute string) error {
	logger.Infof("nextRoute: %v\n", nextRoute)
	switch nextRoute {
	case "api":
		c.Next(api.Handler)
	default:
		return c.NotFound()
	}
	return nil
}
