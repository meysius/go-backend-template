package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := NewApp()
	defer app.pool.Close()

	r := gin.Default()
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

	r.Run(":" + app.cfg.PORT)
}
