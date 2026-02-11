package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string
}

func newErrorResponse(ctx *gin.Context, handlerName string, statusCode int, message string) {
	logrus.WithFields(logrus.Fields{
		"error":   message,
		"handler": handlerName,
		"path":    ctx.Request.URL.Path,
		"method":  ctx.Request.Method,
	}).Warn("handler error")
	ctx.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
