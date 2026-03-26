package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"go-starter-template/controllers"
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

func (a *App) Mount(r *gin.Engine, path string) {
	api := r.Group(path)
	a.usersController.RegisterRoutes(api.Group("/users"))
	a.productsController.RegisterRoutes(api.Group("/products"))
}
