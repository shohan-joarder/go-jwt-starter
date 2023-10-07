package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shohan-joarder/go-jwt-starter/controllers"
	"github.com/shohan-joarder/go-jwt-starter/middlewares"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	// load env here
	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8081"
	}

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Welcome to api routing"})
	})

	api.POST("/register",controllers.Register)
	api.POST("/login",controllers.Login)

	autCheck := api.Use(middlewares.AuthMiddleware())
	autCheck.GET("/user", func(c *gin.Context) {
		fmt.Println(c);
		user, _ := c.Get("email")
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Hello, %s!", user)})
	})

	router.Run(":" + port)

}