package core

import (
	"github.com/goplume/core/health"
	utils2 "github.com/goplume/core/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BrowseLog(ctx *gin.Context) {
	file, isSet := ctx.GetQuery("f")
	if !isSet || file == "" {
		file = "service.log"
	}
	str, err := utils2.ReadLog(file)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
	}
	ctx.String(http.StatusOK, str)
}


//func commonHandler(RouterGroup *gin.RouterGroup, healthFunction health.HealthFunc, middleware ...gin.HandlerFunc) {
//	RouterGroup.GET(rest_api.Endpoint_Log, BrowseLog)
//	RouterGroup.GET(rest_api.Endpoint_Api, RedirectToApi)
//	RouterGroup.GET(rest_api.Endpoint_Api_Any, ginSwagger.WrapHandler(swaggerFiles.Handler))
//	RouterGroup.GET(rest_api.Endpoint_Root, RedirectToApi)
//	if healthFunction != nil {
//		RouterGroup.GET(rest_api.Endpoint_Health, HealthHandler(healthFunction))
//	}
//	RouterGroup.Use(middleware...)
//}

func RedirectToApiHandler(RouterGroup *gin.RouterGroup) {

}

// Health godoc
// @_Summary Show health service status
// @_Produce  json
// @_Success 200 {string} nil ""
// @_Failure 400 {object} fault.TypedErrorStr ""
// @_Failure 404 {object} fault.TypedErrorStr ""
// @_Failure 500 {object} fault.TypedErrorStr ""
// @_Router /health [get]
func HealthHandler(healthFunction health.HealthFunc) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		healthResult := healthFunction()
		if healthResult.IsUp() {
			ctx.JSON(http.StatusOK, healthResult)
		} else if healthResult.IsDown() {
			ctx.JSON(521, healthResult)
		} else if healthResult.IsOutOfService() {
			ctx.JSON(http.StatusOK, healthResult)
		} else if healthResult.IsUnknown() {
			ctx.JSON(http.StatusInternalServerError, healthResult)
		}
	}

}
