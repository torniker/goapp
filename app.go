package wrap

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/torniker/wrap/logger"
	"github.com/torniker/wrap/request"
	"github.com/torniker/wrap/response"
)

// Progliication environments
const (
	Production  string = "Production"
	Testing     string = "Testing"
	Development string = "Development"
)

// Prog contains application data
type Prog struct {
	Env            string
	Server         *http.Server
	Store          map[string]interface{}
	DefaultHandler HandlerFunc
}

var p Prog

// New creates Prog instance and sets default handler function
func New() *Prog {
	p = Prog{
		Env:    Development,
		Server: new(http.Server),
		Store:  make(map[string]interface{}),
		DefaultHandler: func(c *Ctx) error {
			return c.NotFound()
		},
	}
	return &p
}

// Instance returns latest created instance of app
func Instance() *Prog {
	return &p
}

// StartHTTP starts web server
func (p *Prog) StartHTTP(address string) error {
	p.Server.Addr = address
	p.Server.Handler = p
	p.Server.ReadTimeout = 5 * time.Second
	p.Server.WriteTimeout = 10 * time.Second
	return p.Server.ListenAndServe()
}

// StartTLS starts web server with https support
func (p *Prog) StartTLS(address, securedAddr string) error {
	p.Server.Addr = securedAddr
	p.Server.Handler = p
	p.Server.ReadTimeout = 5 * time.Second
	p.Server.WriteTimeout = 10 * time.Second
	go http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://localhost"+securedAddr+r.RequestURI, http.StatusMovedPermanently)
	}))
	return p.Server.ListenAndServeTLS("certs/server.crt", "certs/server.key")
}

// NewCtx returns pointer to Ctx
func (p *Prog) NewCtx(req request.Request, resp response.Response) *Ctx {
	return &Ctx{
		Prog:     p,
		Request:  req,
		Response: resp,
	}
}

func (p *Prog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := p.NewCtx(request.NewHTTP(r), response.NewHTTP(w))
	logger.Infof("---> method: %v, path: %v, query: %v", ctx.Request.Action().String(), ctx.Request.Path(), ctx.Request.Flags())
	err := p.DefaultHandler(ctx)
	if err != nil {
		ctx.Error(err)
	}
}

// StartCLI starts cli server
func (p *Prog) StartCLI() error {
	buf := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	fmt.Println("--------------------------------------------------------")
	fmt.Println("Please provide command in following format:")
	fmt.Println("--> post command/path?flag1=foo&flag2=bar {\"input\":\"data\"}")
	fmt.Println("--------------------------------------------------------")
	return p.ListenCLI(buf)
}

// ListenCLI listens to buffer input
func (p *Prog) ListenCLI(buf *bufio.ReadWriter) error {
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
		var addr string
		var action request.Action
		var input io.Reader
		if len(parts) == 2 {
			action = request.READ
			addr = "http://app.cli/" + strings.TrimSpace(parts[0])
			input = strings.NewReader(parts[1])
		} else if len(parts) == 3 {
			action = request.NewActionFromString(parts[0])
			if !action.IsValid() {
				fmt.Println("invalid action")
				continue
			}
			addr = "http://app.cli/" + strings.TrimSpace(parts[1])
			input = strings.NewReader(parts[1])
		} else {
			fmt.Println("could not understand query")
			continue
		}
		u, err := url.Parse(addr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ctx := p.NewCtx(request.NewCLI(action, u, input), response.NewCLI())
		logger.Infof("---> action: %v, path: %v, query: %v", ctx.Request.Action().String(), ctx.Request.Path(), ctx.Request.Flags())
		err = p.DefaultHandler(ctx)
		if err != nil {
			ctx.Error(err)
		}
	}
}

// IsDevelopment is true if the env is development
func (p *Prog) IsDevelopment() bool {
	return p.Env == Development
}

// IsTesting is true if the env is testing
func (p *Prog) IsTesting() bool {
	return p.Env == Testing
}

// IsProduction is true if the env is production
func (p *Prog) IsProduction() bool {
	return p.Env == Production
}

// Call helps to create new request
func (p *Prog) Call() *Requester {
	return &Requester{
		prog: p,
		req:  &request.Req{},
	}
}

// Requester wraps sub requests
type Requester struct {
	prog *Prog
	req  *request.Req
	err  error
}

// Bind tries to call the request and bind the data from the response
func (r *Requester) Bind(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	subCtx := r.prog.NewCtx(r.req, response.NewResponse())
	err := r.prog.DefaultHandler(subCtx)
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
