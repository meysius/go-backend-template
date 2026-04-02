package main

import (
	"context"
	"log"
	"log/slog"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"go-starter-template/controllers"
	"go-starter-template/docs"
	"go-starter-template/domain/identity"
	"go-starter-template/domain/ordering"
)

type App struct {
	cfg                *Config
	logger             *slog.Logger
	pool               *pgxpool.Pool
	identityService    *identity.IdentityService
	usersController    *controllers.UsersController
	orderingService    *ordering.OrderingService
	productsController *controllers.ProductsController
}

func NewApp() *App {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg, err := Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := NewLogger(cfg)
	logger.Info("starting", "env", cfg.ENV)
	logger.Info("Configuration", "config", cfg)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	identityRepo := identity.NewIdentityRepo(pool)
	identityService := identity.NewIdentityService(identityRepo)

	orderingRepo := ordering.NewOrderingRepo(pool)
	orderingService := ordering.NewOrderingService(orderingRepo)

	return &App{
		cfg:                cfg,
		logger:             logger,
		pool:               pool,
		identityService:    identityService,
		usersController:    controllers.NewUsersController(identityService),
		orderingService:    orderingService,
		productsController: controllers.NewProductsController(orderingService),
	}
}

const scalarHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>Go Starter API</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/openapi.json"
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

func (a *App) Mount(r *gin.Engine, path string) {
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(scalarHTML))
	})
	r.GET("/docs/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", []byte(docs.SwaggerInfo.ReadDoc()))
	})

	r.GET("/health", func(c *gin.Context) {
		if err := a.pool.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group(path)
	a.usersController.RegisterRoutes(api.Group("/users"))
	a.productsController.RegisterRoutes(api.Group("/products"))
}
