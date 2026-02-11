package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func (h *GoodsHandlers) UserIdentity(ctx *gin.Context) {
	nameHandler := "UserIdentity"
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "token is empty")
		return
	}

	token := headerParts[1]

	response, err := h.authClient.ValidateToken(ctx, token)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}
	if !response.Valid {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := strconv.Atoi(response.UserId)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "invalid user id")
		return
	}

	ctx.Set("userID", id)
	ctx.Set("userRole", response.Role)
	ctx.Set("userName", response.Us)
}
