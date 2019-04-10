package main

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/should"
	"github.com/torniker/goapp/web"
)

func main() {
	a, err := app.New(web.Handler)
	if err != nil {
		panic(err)
	}
	a.Start(":8989")
}

func home(c *app.Ctx, nextRoute string) error {
	logger.Infof("nextRoute: %v\n", nextRoute)
	c.Next(api)
	return nil
}

func api(c *app.Ctx, nextRoute string) error {
	logger.Infof("nextRoute: %v\n", nextRoute)
	should.Method(c, should.POST)
	c.Next(auth)
	return nil
}

func auth(c *app.Ctx, nextRoute string) error {
	logger.Infof("nextRoute: %v\n", nextRoute)
	return nil
}
