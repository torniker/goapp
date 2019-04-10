package web

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/web/api"
)

func Handler(c *app.Ctx, nextSegment string) error {
	if nextSegment == "api" {
		c.Next(api.Handler)
		return nil
	}
	return c.NotFound()
}
