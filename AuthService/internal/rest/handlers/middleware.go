package handlers

import (
	"errors"
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
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, "UserIdentity", http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(ctx, "UserIdentity", http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	userID, role, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, "UserIdentity", http.StatusUnauthorized, err)
	}
	ctx.Set(userIDCtx, userID)
	ctx.Set(userRoleCtx, role)
}
