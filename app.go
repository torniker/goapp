package app

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/torniker/goapp/logger"
	"github.com/torniker/goapp/request"
	"github.com/torniker/goapp/response"
)

// Appliication environments
const (
	Production  string = "Production"
	Testing     string = "Testing"
	Development string = "Development"
)

// App contains application data
type App struct {
	Env            string
	Server         *http.Server
	Store          map[string]interface{}
	DefaultHandler HandlerFunc
}

var a App

// New creates App instance and sets default handler function
func New() *App {
	a = App{
		Env:    Development,
		Server: new(http.Server),
		Store:  make(map[string]interface{}),
		DefaultHandler: func(c *Ctx) error {
			return c.NotFound()
		},
	}
	return &a
}

// Instance returns latest created instance of app
func Instance() *App {
	return &a
}

// StartHTTP starts web server
func (a *App) StartHTTP(address string) error {
	a.Server.Addr = address
	a.Server.Handler = a
	a.Server.ReadTimeout = 5 * time.Second
	a.Server.WriteTimeout = 10 * time.Second
	return a.Server.ListenAndServe()
}

// StartTLS starts web server with https support
func (a *App) StartTLS(address, securedAddr string) error {
	a.Server.Addr = securedAddr
	a.Server.Handler = a
	a.Server.ReadTimeout = 5 * time.Second
	a.Server.WriteTimeout = 10 * time.Second
	go http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://localhost"+securedAddr+r.RequestURI, http.StatusMovedPermanently)
	}))
	return a.Server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
}

// NewCtx returns pointer to Ctx
func (a *App) NewCtx(req request.Request, resp response.Response) *Ctx {
	return &Ctx{
		App:      a,
		Request:  req,
		Response: resp,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := a.NewCtx(request.NewHTTP(r), response.NewHTTP(w))
	logger.Infof("---> method: %v, path: %v, query: %v", ctx.Request.Action().String(), ctx.Request.Path(), ctx.Request.Flags())
	err := a.DefaultHandler(ctx)
	if err != nil {
		ctx.Error(err)
	}
}

// StartCLI starts cli server
func (a *App) StartCLI() error {
	buf := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	fmt.Println("--------------------------------------------------------")
	fmt.Println("Please provide command in following format:")
	fmt.Println("--> post command/path?flag1=foo&flag2=bar {\"input\":\"data\"}")
	fmt.Println("--------------------------------------------------------")
	return a.ListenCLI(buf)
}

// ListenCLI listens to buffer input
func (a *App) ListenCLI(buf *bufio.ReadWriter) error {
	for {
		fmt.Print("--> ")
		commandBytes, err := buf.ReadBytes('\n')
		if err != nil {
			return err
		}
		command := string(commandBytes)
		if strings.ToLower(strings.TrimSpace(command)) == "exit" {
			return nil
		}
		parts := strings.SplitN(command, " ", 3)
		var p string
		var action request.Action
		var input io.Reader
		if len(parts) == 2 {
			action = request.READ
			p = "http://app.cli/" + strings.TrimSpace(parts[0])
			input = strings.NewReader(parts[1])
		} else if len(parts) == 3 {
			action = request.NewActionFromString(parts[0])
			if !action.IsValid() {
				fmt.Println("invalid action")
				continue
			}
			p = "http://app.cli/" + strings.TrimSpace(parts[1])
			input = strings.NewReader(parts[1])
		} else {
			fmt.Println("could not understand query")
			continue
		}
		u, err := url.Parse(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ctx := a.NewCtx(request.NewCLI(action, u, input), response.NewCLI())
		logger.Infof("---> action: %v, path: %v, query: %v", ctx.Request.Action().String(), ctx.Request.Path(), ctx.Request.Flags())
		err = a.DefaultHandler(ctx)
		if err != nil {
			ctx.Error(err)
		}
	}
}

// IsDevelopment is true if the env is development
func (a *App) IsDevelopment() bool {
	return a.Env == Development
}

// IsTesting is true if the env is testing
func (a *App) IsTesting() bool {
	return a.Env == Testing
}

// IsProduction is true if the env is production
func (a *App) IsProduction() bool {
	return a.Env == Production
}

// Call helps to create new request
func (a *App) Call() *Requester {
	return &Requester{
		app: a,
		req: &request.Req{},
	}
}

// Requester wraps sub requests
type Requester struct {
	app *App
	req *request.Req
	err error
}

// Bind tries to call the request and bind the data from the response
func (r *Requester) Bind(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	subCtx := r.app.NewCtx(r.req, response.NewResponse())
	err := r.app.DefaultHandler(subCtx)
	if err != nil {
		return err
	}
	v = subCtx.Response.Output()
	return nil
}

// Path sets command name for the request
func (r *Requester) path(command string) *Requester {
	u, err := url.Parse(command)
	if err != nil {
		r.err = err
		return r
	}
	r.req.SetPath(u)
	return r
}

// Flags sets flags for the request
func (r *Requester) Flags(flags map[string][]string) *Requester {
	r.req.SetFlags(flags)
	return r
}

// Input sets input data for request
func (r *Requester) Input(v interface{}) *Requester {
	r.req.SetData(v)
	return r
}

// Create command
func (r *Requester) Create(command string) *Requester {
	r.req.SetAction(request.CREATE)
	r.path(command)
	return r
}

// Read command
func (r *Requester) Read(command string) *Requester {
	r.req.SetAction(request.READ)
	r.path(command)
	return r
}

// Update command
func (r *Requester) Update(command string) *Requester {
	r.req.SetAction(request.UPDATE)
	r.path(command)
	return r
}

// Delete command
func (r *Requester) Delete(command string) *Requester {
	r.req.SetAction(request.DELETE)
	r.path(command)
	return r
}
