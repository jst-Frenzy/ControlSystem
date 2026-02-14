package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jst-Frenzy/ControlSystem/AuthService/internal/AuthService"
	"net/http"
)

type AuthHandler struct {
	service AuthService.AuthService
}

func NewAuthHandler(service AuthService.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {
	nameHandler := "SignUp"
	var user AuthService.UserSignUp
	if err := ctx.ShouldBind(&user); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.service.SignUp(user)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *AuthHandler) SignIn(ctx *gin.Context) {
	nameHandler := "SignIn"
	var user AuthService.UserSignIn
	if err := ctx.ShouldBind(&user); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.service.SignIn(user)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access token":  tokens.AccessToken,
		"refresh token": tokens.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(ctx *gin.Context) {
	nameHandler := "Refresh"
	var refreshToken AuthService.RefreshTokenRequest
	if err := ctx.ShouldBind(&refreshToken); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	newAccessToken, err := h.service.RefreshTokens(refreshToken.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access token": newAccessToken,
	})
}

func (h *AuthHandler) ChangeRole(ctx *gin.Context) {
	nameHandler := "ChangeRole"
	type data struct {
		User    AuthService.UserSignIn `json:"user"`
		NewRole string                 `json:"newRole"`
		Id      int                    `json:"id"`
	}
	var d data
	if err := ctx.ShouldBind(&d); err != nil {
		newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.ChangeRole(d.User, d.Id, d.NewRole)
	if err != nil {
		if errors.Is(err, errors.New("not enough rights")) {
			newErrorResponse(ctx, nameHandler, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(ctx, nameHandler, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
