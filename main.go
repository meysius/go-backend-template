package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-starter-template/app"
	"go-starter-template/middleware"
)

// @title           Go Starter API
// @version         1.0
// @description     A Go web API starter template using Gin, PostgreSQL, sqlc, and goose.
// @host            localhost:3000
// @BasePath        /api
func main() {
	a := app.NewApp()

	if a.Cfg.ENV == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(a.Logger))
	a.Mount(r, "/api")

	r.GET("/test", func(c *gin.Context) {
		a.Logger.Info("handling test request")
		name := fmt.Sprintf("user-%d", rand.Intn(100000))
		user, err := a.IdentityService.CreateUser(c.Request.Context(), name, name+"@example.com")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	srv := &http.Server{
		Addr:    ":" + a.Cfg.PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	a.Logger.Info("Gin server started", "env", a.Cfg.ENV, "address", fmt.Sprintf("http://localhost:%s", a.Cfg.PORT))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.Logger.Error("server forced to shutdown", "error", err)
	}

	a.Pool.Close()
	a.Logger.Info("server stopped")
}
