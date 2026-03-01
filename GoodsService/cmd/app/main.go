package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/dataBase"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/gRPC/client"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/gRPC/server"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/rest/handlers"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logrus.New()
	dataBase.InitMongo()

	goodsMongoRepo := GoodService.NewGoodsMongoRepo(dataBase.MongoDB)
	goodsService := GoodService.NewGoodService(goodsMongoRepo)
	authClientGRPC, err := client.NewAuthClient(os.Getenv("ADDRESS_GRPC_SERVER"))
	if err != nil {
		logrus.Fatal("Cant start grpc client")
	}

	grpcServer := server.NewGRPCServer(server.Deps{
		GoodsService: goodsService,
		Logger:       logger,
	})

	go func() {
		port := 50052

		if errStart := grpcServer.Start(port); err != nil {
			logger.WithError(errStart).Fatal("grpc server failed")
		}
	}()

	goodsHandler := handlers.NewGoodsHandlers(goodsService, authClientGRPC)

	router := gin.Default()

	api := router.Group("/api/goods")
	{
		api.GET("/catalog", goodsHandler.GetGoods)

		itemGroup := api.Group("/item")
		itemGroup.Use(goodsHandler.UserIdentity)
		{
			itemGroup.POST("/", goodsHandler.AddItem)
			itemGroup.DELETE("/:id", goodsHandler.DeleteItem)
			itemGroup.PUT("/", goodsHandler.UpdateItem)
		}
	}

	go func() {
		if errRun := router.Run(":8081"); errRun != nil {
			logrus.WithError(errRun).Fatalf("REST server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	grpcServer.Stop()

	time.Sleep(2 * time.Second)
}
