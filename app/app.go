package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-starter-template/controllers"
	"go-starter-template/docs"
	"go-starter-template/domain/identity"
	"go-starter-template/views"
)

type App struct {
	Cfg                *Config
	Logger             *slog.Logger
	Pool               *pgxpool.Pool
	IdentityService *identity.IdentityService
	usersController *controllers.UsersController
}

func NewApp() *App {
	cfg := loadConfig()
	logger := newLogger(cfg.ENV)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		logger.Error("Failed to create connection pool", "error", err)
		os.Exit(1)
	}
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	identityRepo := identity.NewIdentityRepo(pool)
	identityService := identity.NewIdentityService(identityRepo)

	return &App{
		Cfg:             cfg,
		Logger:          logger,
		Pool:            pool,
		IdentityService: identityService,
		usersController: controllers.NewUsersController(identityService),
	}
}

func (a *App) Mount(r *gin.Engine, path string) {
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", views.MustRead("welcome.html"))
	})

	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", views.MustRead("scalar.html"))
	})
	r.GET("/docs/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	r.GET("/health", func(c *gin.Context) {
		if err := a.Pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group(path)
	a.usersController.RegisterRoutes(api.Group("/users"))
}
