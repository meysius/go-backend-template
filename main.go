package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-starter-template/middleware"
)

// @title           Go Starter API
// @version         1.0
// @description     A Go web API starter template using Gin, PostgreSQL, sqlc, and goose.
// @host            localhost:3000
// @BasePath        /api
func main() {
	app := NewApp()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger())
	app.Mount(r, "/api")

	r.GET("/test", func(c *gin.Context) {
		name := fmt.Sprintf("user-%d", rand.Intn(100000))
		user, err := app.identityService.CreateUser(name, name+"@example.com")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	srv := &http.Server{
		Addr:    ":" + app.cfg.PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("server started", "address", fmt.Sprintf("http://localhost:%s", app.cfg.PORT))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	app.pool.Close()
	slog.Info("server stopped")
}
