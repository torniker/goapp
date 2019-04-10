package should

import "github.com/torniker/goapp/app"

// http request methods
const (
	GET int = 1 << iota
	POST
	PUT
	DELETE
)

// Method checks if the method is
// should.Method(ctx, should.GET | should.POST)
func Method(c *app.Ctx, m int) {
	if m&methodFromString(c.Method()) == 0 {
		c.NotFound()
	}
}

func methodFromString(m string) int {
	switch m {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	}
	return 0
}
