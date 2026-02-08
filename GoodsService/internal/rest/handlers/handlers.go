package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"net/http"
)

type GoodsHandlers struct {
	serv GoodService.GoodService
}

func NewGoodsHandlers(serv GoodService.GoodService) *GoodsHandlers {
	return &GoodsHandlers{serv: serv}
}

func (gh *GoodsHandlers) GetGoods(ctx *gin.Context) {
	nameHandler := "GetGoods"
	items, err := gh.serv.GetGoods()
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err)
	}

	resp := make(map[string]string)

	for _, i := range items {
		resp[i.Name] = i.Description
	}

	ctx.JSON(http.StatusOK, resp)
}

func (gh *GoodsHandlers) AddItem() {

}
