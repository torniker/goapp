package app

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/olivere/elastic"
	"github.com/torniker/goapp/app/logger"
	"github.com/torniker/goapp/app/request"
	"github.com/torniker/goapp/app/response"
)

type Environment string

const (
	Production  Environment = "Production"
	Testing     Environment = "Testing"
	Development Environment = "Development"
)

// App contains application data
type App struct {
	Server         *http.Server
	DefaultHandler HandlerFunc

	Env      Environment
	postgres *sqlx.DB
	elastic  *elastic.Client
	redis    *redis.Client
}

// New creates App instance and sets default handler function
func New() *App {
	return &App{
		Server: new(http.Server),
		DefaultHandler: func(c *Ctx, nextRoute string) error {
			return c.NotFound()
		},
	}
}

func defaultHandler(c *Ctx, nextRoute string) error {
	return c.NotFound()
}

// Start starts web server
func (a *App) Start(address string) error {
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
func (a *App) NewCtx() *Ctx {
	return &Ctx{
		App: a,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := a.NewCtx()
	ctx.request = request.NewHTTP(r)
	ctx.response = response.NewHTTP(w)
	ctx.path = &path{
		path:     r.URL.Path,
		segments: strings.Split(r.URL.Path, "/"),
		index:    0,
	}
	logger.Infof("---> method: %v, path: %v, query: %v", r.Method, r.URL.Path, r.URL.Query())
	err := a.DefaultHandler(ctx, ctx.path.Next())
	if err != nil {
		ctx.Error(err)
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

func (a *App) InitPG(addr string) error {
	_ = pq.Efatal
	db, err := sqlx.Connect("postgres", addr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	a.postgres = db
	return nil
}

// PG returns postgres db object
func (a *App) PG() *sqlx.DB {
	return a.postgres
}

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
