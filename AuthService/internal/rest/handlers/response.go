package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string
}

func newErrorResponse(ctx *gin.Context, handlerName string, statusCode int, err error) {
	logrus.WithFields(logrus.Fields{
		"error":   err.Error(),
		"handler": handlerName,
		"path":    ctx.Request.URL.Path,
		"method":  ctx.Request.Method,
	}).Warn("handler error")
	ctx.AbortWithStatusJSON(statusCode, errorResponse{Message: err.Error()})
}
