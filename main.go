package main

import (
	"github.com/asimbera/pokket/controllers"
	"github.com/asimbera/pokket/models"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func setupRouter() *gin.Engine {
	if err := models.SetupDatabase(); err != nil {
		os.Exit(1) // Terminate the process
	}

	r := gin.Default()

	api := r.Group("/api") // /api
	v1 := api.Group("/v1") // /api/v1
	{
		v1.GET("/status", controllers.HealthController) // /api/v1/status
		auth := v1.Group("/auth")                       // /api/v1/auth
		{
			auth.POST("/login", controllers.LoginController)   // /api/v1/auth/login
			auth.POST("/signup", controllers.SignupController) // /api/v1/auth/signup
		}
		secure := v1.Group("/secure")
		secure.Use(controllers.AuthMiddleware())
		{
			secure.GET("/me")
		}
	}

	return r
}

func main() {
	r := setupRouter()

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
