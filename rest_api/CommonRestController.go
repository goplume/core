package rest_api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonRestController interface {
	GetList(ctx *gin.Context)
	MassUpdate(ctx *gin.Context)
	MassDelete(ctx *gin.Context)
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func FunctionalNotImplemented(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"status": "Must be implemented in next version"})
}

func FunctionalNotAllowed(ctx *gin.Context) {
	ctx.JSON(http.StatusMethodNotAllowed, gin.H{"status": "MethodNotAllowed"})
}
