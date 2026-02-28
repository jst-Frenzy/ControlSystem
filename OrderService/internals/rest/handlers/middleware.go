package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func (h *OrderHandler) UserIdentity(ctx *gin.Context) {
	handlerName := "UserIdentity"
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, handlerName, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(ctx, handlerName, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(ctx, handlerName, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(ctx, handlerName, http.StatusUnauthorized, "token is empty")
		return
	}

	token := headerParts[1]

	resp, err := h.authClient.ValidateToken(ctx, token)
	if err != nil {
		newErrorResponse(ctx, handlerName, http.StatusBadRequest, err.Error())
		return
	}
	if !resp.Valid {
		newErrorResponse(ctx, handlerName, http.StatusBadRequest, err.Error())
		return
	}

	ctx.Set("CartID", resp.CartId)
}
