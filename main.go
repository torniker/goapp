package main

import (
	"github.com/torniker/goapp/app"
	"github.com/torniker/goapp/web"
)

func main() {
	a, err := app.New(web.Handler)
	if err != nil {
		panic(err)
	}
	a.Start(":8989")
}
