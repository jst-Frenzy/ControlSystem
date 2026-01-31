package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello World!")
}
