package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/dataBase"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/gRPC"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/rest/handlers"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	dataBase.InitMongo()

	goodsMongoRepo := GoodService.NewGoodsMongoRepo(dataBase.MongoDB)
	goodsService := GoodService.NewGoodService(goodsMongoRepo)
	authClientGRPC, err := gRPC.NewAuthClient(os.Getenv("ADDRESS_GRPC_SERVER"))
	if err != nil {
		logrus.Fatal("Cant start grpc client")
	}

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

	if errRun := router.Run(":8080"); errRun != nil {
		logrus.WithError(errRun).Fatalf("REST server failed")
	}
}
