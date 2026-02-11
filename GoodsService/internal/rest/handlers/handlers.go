package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/GoodService"
	"github.com/jst-Frenzy/ControlSystem/GoodsService/internal/gRPC"
	"net/http"
)

type GoodsHandlers struct {
	serv       GoodService.GoodService
	authClient *gRPC.AuthClient
}

func NewGoodsHandlers(serv GoodService.GoodService, authClient *gRPC.AuthClient) *GoodsHandlers {
	return &GoodsHandlers{
		serv:       serv,
		authClient: authClient,
	}
}

func (h *GoodsHandlers) GetGoods(ctx *gin.Context) {
	nameHandler := "GetGoods"
	items, err := h.serv.GetGoods()
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make(map[string]string)

	for _, i := range items {
		resp[i.Name] = i.Description
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *GoodsHandlers) AddItem(ctx *gin.Context) {
	nameHandler := "AddItem"
	role := ctx.MustGet("userRole")

	if !(role == "seller") {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "not enough rights")
		return
	}

	var i GoodService.Item
	if err := ctx.ShouldBind(&i); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.serv.AddItem(i)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *GoodsHandlers) DeleteItem(ctx *gin.Context) {
	nameHandler := "DeleteItem"
	role := ctx.MustGet("userRole")

	if role != "seller" {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "not enough rights")
		return
	}

	itemID, ok := ctx.GetQuery("id")
	if !ok {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "no id for delete item")
	}

	err := h.serv.DeleteItem(itemID)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *GoodsHandlers) UpdateItem(ctx *gin.Context) {
	nameHandler := "UpdateItem"
	role := ctx.MustGet("userRole")

	if role != "seller" {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "not enough rights")
		return
	}

	var i GoodService.Item
	if err := ctx.ShouldBind(&i); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	respItem, err := h.serv.UpdateItem(i)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, respItem)
}
