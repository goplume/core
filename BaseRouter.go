package core

import "github.com/gin-gonic/gin"

type BaseRouter struct {
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
	Filters     []gin.HandlerFunc
}

func (this BaseRouter) GetEngine() *gin.Engine {
	return this.Router
}

func (this *BaseRouter) GetRouterGroup() *gin.RouterGroup {
	return this.RouterGroup
}

func (this *BaseRouter) Get(relativePath string, handlers ...gin.HandlerFunc) {
	this.RouterGroup.GET(relativePath, handlers...)
}

func (this *BaseRouter) AddFilter(handlers ...gin.HandlerFunc) {
	this.RouterGroup.Use(handlers...)
}
