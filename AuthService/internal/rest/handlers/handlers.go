package handlers

import (
	"AuthService/internal/AuthService"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	service AuthService.AuthService
}

func NewAuthHandler(service AuthService.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	var user AuthService.UserSignUp
	if err := ctx.ShouldBind(&user); err != nil {
		newErrorResponse(ctx, "SignUp", http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.service.SignUp(user)
	if err != nil {
		newErrorResponse(ctx, "SignUP", http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var user AuthService.UserSignIn
	if err := ctx.ShouldBind(&user); err != nil {
		newErrorResponse(ctx, "SignIn", http.StatusBadRequest, err)
		return
	}

	tokens, err := h.service.SignIn(user)
	if err != nil {
		newErrorResponse(ctx, "SignIn", http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	var refreshToken AuthService.RefreshTokenRequest
	if err := ctx.ShouldBind(&refreshToken); err != nil {
		newErrorResponse(ctx, "Refresh", http.StatusBadRequest, err)
		return
	}

	newAccessToken, err := h.service.RefreshTokens(refreshToken.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, "Refresh", http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access token": newAccessToken,
	})
}
