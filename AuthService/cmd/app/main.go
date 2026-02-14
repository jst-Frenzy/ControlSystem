package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/dataBase"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/gRPC"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/rest/handlers"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logrus.New()
	dataBase.InitRedis()
	dataBase.InitPostgres()

	authPostgresRepo := AuthService.NewAuthPostgresRepo(dataBase.PostgresDB)
	authRedisRepo := AuthService.NewAuthRedisRepo(dataBase.RedisDB)
	tokenManager := AuthService.NewManager(os.Getenv("SIGNING_KEY"))
	authService := AuthService.NewAuthService(authPostgresRepo, authRedisRepo, tokenManager)

	grpcServer := gRPC.NewGRPCServer(gRPC.Deps{
		Logger:      logger,
		AuthService: authService,
	})

	go func() {
		port := 50051

		if err := grpcServer.StartGRPC(port); err != nil {
			logger.WithError(err).Fatal("gRPC server failed")
		}
	}()

	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/signin", authHandler.SignIn)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/changeRole", authHandler.ChangeRole)
		}
	}

	go func() {
		if err := router.Run(":8080"); err != nil {
			logrus.WithError(err).Fatalf("REST server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.Stop()

	time.Sleep(2 * time.Second)
}
