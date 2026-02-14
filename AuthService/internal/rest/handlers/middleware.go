package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userIDCtx           = "userID"
	userRoleCtx         = "userRole"
)

func (h *AuthHandler) UserIdentity(ctx *gin.Context) {
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

	info, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusUnauthorized, err.Error())
		return
	}
	ctx.Set(userIDCtx, info.ID)
	ctx.Set(userRoleCtx, info.Role)
}
