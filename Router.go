package core

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Get(relativePath string, handlers ...gin.HandlerFunc)
	GetEngine() *gin.Engine
	GetRouterGroup() *gin.RouterGroup
	AddFilter(handlers ...gin.HandlerFunc)
	InitRouter()
}
