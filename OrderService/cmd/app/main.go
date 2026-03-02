package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/dataBase"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/gRPC/client"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/orderService"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/rest/handlers"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	dataBase.InitPostgres()

	orderPostgresRepo := orderService.NewOrderPostgresRep(dataBase.PostgresDB)

	goodsClient, err := client.NewGoodsClient(os.Getenv("ADDRESS_GRPC_GOODS"))
	if err != nil {
		logrus.Fatal("can't start goods client")
	}

	orderServ := orderService.NewOrderService(orderPostgresRepo, goodsClient)

	authClient, err := client.NewAuthClient(os.Getenv("ADDRESS_GRPC_AUTH"))
	if err != nil {
		logrus.Fatal("can't start auth client")
	}

	orderHandler := handlers.NewOrderHandler(orderServ, authClient)

	router := gin.Default()

	api := router.Group("/api/orders")
	api.Use(orderHandler.UserIdentity)
	{
		cartGroup := api.Group("/cart")
		{
			cartGroup.GET("/", orderHandler.GetCart)
			cartGroup.POST("/", orderHandler.AddToCart)
			cartGroup.DELETE("/:id", orderHandler.DeleteFromCart)
		}
	}

	if errRun := router.Run(":8082"); errRun != nil {
		logrus.WithError(err).Fatalf("can't start order server")
	}
}
