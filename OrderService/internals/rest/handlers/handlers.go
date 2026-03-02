package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/gRPC/client"
	"github.com/jst-Frenzy/ControlSystem/OrderService/internals/orderService"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	serv       orderService.OrderService
	authClient client.AuthClient
}

func NewOrderHandler(serv orderService.OrderService, authClient client.AuthClient) *OrderHandler {
	return &OrderHandler{
		serv:       serv,
		authClient: authClient,
	}
}

func (h *OrderHandler) AddToCart(ctx *gin.Context) {
	nameHandler := "AddToCart"
	cartIDstr := ctx.MustGet("CartID").(string)
	cartID, _ := strconv.Atoi(cartIDstr)

	var i orderService.CartItem
	if err := ctx.ShouldBind(&i); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	i.CartID = cartID
	id, err := h.serv.AddToCart(i)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *OrderHandler) GetCart(ctx *gin.Context) {
	nameHandler := "GetCart"
	cartIDstr := ctx.MustGet("CartID").(string)
	cartID, _ := strconv.Atoi(cartIDstr)

	cart, totalPrice, err := h.serv.GetCart(cartID, ctx)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
	}

	type ItemStruct struct {
		ProductID string
		Quantity  int
		Price     float64
	}

	resp := make(map[string]interface{})

	for _, i := range cart {
		resp[i.Name] = ItemStruct{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
			Price:     i.Price,
		}
	}

	resp["total price"] = totalPrice

	ctx.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) DeleteFromCart(ctx *gin.Context) {
	nameHandler := "DeleteFromCart"
	cartIDstr := ctx.MustGet("CartID").(string)
	cartID, _ := strconv.Atoi(cartIDstr)

	itemID := ctx.Param("id")

	err := h.serv.RemoveFromCart(cartID, itemID)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
