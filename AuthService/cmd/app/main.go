package main

import (
	"AuthService/internal/AuthService"
	"AuthService/internal/dataBase"
	"AuthService/internal/rest/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	dataBase.InitRedis()
	dataBase.InitPostgres()

	authPostgresRepo := AuthService.NewAuthPostgresRepo(dataBase.PostgresDB)
	authRedisRepo := AuthService.NewAuthRedisRepo(dataBase.RedisDB)
	tokenManager := AuthService.NewManager(os.Getenv("SIGNING_KEY"))
	authService := AuthService.NewAuthService(authPostgresRepo, authRedisRepo, tokenManager)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/signin", authHandler.SignIn)
			auth.POST("/refresh", authHandler.Refresh)
		}
	}

	err := router.Run(":8080")
	if err != nil {
		logrus.WithError(err).Fatalf("Can not to start Auth server")
	}
}
