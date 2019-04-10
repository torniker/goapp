package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/olivere/elastic"
	"github.com/torniker/goapp/app/logger"
)

type environment string

const (
	Production  environment = "Production"
	Testing     environment = "Testing"
	Development environment = "Development"
)

type Config struct {
	Environment      environment
	ESAddress        string
	RedisAddr        string
	RedisPassword    string
	PostgresAddr     string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
}

// App contains application data
type App struct {
	Server         *http.Server
	DefaultHandler HandlerFunc

	env      environment
	postgres *sqlx.DB
	elastic  *elastic.Client
	redis    *redis.Client
}

// New creates App instance and sets default handler function
func New(f HandlerFunc) (*App, error) {
	a := &App{
		Server:         new(http.Server),
		DefaultHandler: f,
	}
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		return nil, err
	}
	cfg := Config{
		Environment:      environment(os.Getenv("ENV")),
		ESAddress:        os.Getenv("ELASTICSEARCH_ADDRESS"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		PostgresAddr:     os.Getenv("POSTGRES_ADDRESS"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
	}
	logger.Infof("starting instance with config: %#v", cfg)
	a.env = cfg.Environment
	err = a.initPG(cfg)
	if err != nil {
		return nil, err
	}
	err = a.initES(cfg)
	if err != nil {
		return nil, err
	}
	err = a.initRedis(cfg)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Start starts web server
func (a *App) Start(address string) error {
	a.Server.Addr = address
	a.Server.Handler = a
	return a.Server.ListenAndServe()
}

// NewCtx returns pointer to Ctx
func (a *App) NewCtx() *Ctx {
	return &Ctx{
		App: a,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := a.NewCtx()
	ctx.request = r
	ctx.response = w
	ctx.path = path{
		path:     r.URL.Path,
		segments: strings.Split(r.URL.Path, "/"),
		index:    0,
	}
	err := a.DefaultHandler(ctx, ctx.path.Next())
	if err != nil {
		ctx.Error(err)
	}
}

func (a *App) Env() environment {
	return a.env
}

// IsDevelopment is true if the env is development
func (a *App) IsDevelopment() bool {
	return a.env == Development
}

// IsTesting is true if the env is testing
func (a *App) IsTesting() bool {
	return a.env == Testing
}

// IsProduction is true if the env is production
func (a *App) IsProduction() bool {
	return a.env == Production
}

func (a *App) initRedis(cfg Config) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		logger.Warn(err)
		return err
	}
	a.redis = client
	return nil
}

// Redis returns redis connection client
func (a *App) Redis() *redis.Client {
	return a.redis
}

func (a *App) initPG(cfg Config) error {
	_ = pq.Efatal
	addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresAddr, cfg.PostgresDB)
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

func (a *App) initES(cfg Config) error {
	fmt.Println(cfg.ESAddress)
	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s", cfg.ESAddress)),
		elastic.SetSniff(false),
	)
	if err != nil {
		return err
	}
	a.elastic = client
	return nil
}

// ES returns elastic search object
func (a *App) ES() *elastic.Client {
	return a.elastic
}
