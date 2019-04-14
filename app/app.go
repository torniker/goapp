package app

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/app/request"
	"github.com/torniker/goapp/app/response"
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
		CurrentPath: &path{
			path:     req.Path(),
			segments: strings.Split(req.Path(), "/"),
			index:    0,
		},
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := a.NewCtx(request.NewHTTP(r), response.NewHTTP(w))
	logger.Infof("---> method: %v, path: %v, query: %v", ctx.Method(), ctx.Request.Path(), ctx.Request.Query())
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
		var method, p, input string
		if len(parts) == 2 {
			method = "get"
			p = "http://app.cli/" + strings.TrimSpace(parts[0])
			input = parts[1]
		} else if len(parts) == 3 {
			method = parts[0]
			p = "http://app.cli/" + strings.TrimSpace(parts[1])
			input = parts[2]
		} else {
			fmt.Println("could not understand query")
			continue
		}
		u, err := url.Parse(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ctx := a.NewCtx(request.NewCLI(method, u, input), response.NewCLI())
		logger.Infof("---> method: %v, path: %v, query: %v", ctx.Method(), ctx.Request.Path(), ctx.Request.Query())
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

// func (a *App) initRedis(cfg Config) error {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     cfg.RedisAddr,
// 		Password: cfg.RedisPassword,
// 		DB:       0,
// 	})
// 	_, err := client.Ping().Result()
// 	if err != nil {
// 		logger.Warn(err)
// 		return err
// 	}
// 	a.redis = client
// 	return nil
// }

// // Redis returns redis connection client
// func (a *App) Redis() *redis.Client {
// 	return a.redis
// }

// func (a *App) InitPG(addr string) error {
// 	_ = pq.Efatal
// 	db, err := sqlx.Connect("postgres", addr)
// 	if err != nil {
// 		return err
// 	}
// 	err = db.Ping()
// 	if err != nil {
// 		return err
// 	}
// 	a.postgres = db
// 	return nil
// }

// // PG returns postgres db object
// func (a *App) PG() *sqlx.DB {
// 	return a.postgres
// }

// func (a *App) initES(cfg Config) error {
// 	fmt.Println(cfg.ESAddress)
// 	client, err := elastic.NewClient(
// 		elastic.SetURL(fmt.Sprintf("http://%s", cfg.ESAddress)),
// 		elastic.SetSniff(false),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	a.elastic = client
// 	return nil
// }

// // ES returns elastic search object
// func (a *App) ES() *elastic.Client {
// 	return a.elastic
// }
