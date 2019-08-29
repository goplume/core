package rest_api

import (
	"github.com/goplume/core/rest_client"
	"github.com/gin-gonic/gin"
)

type RedirectApiRestController struct {
	BaseRestController
	RestClient *rest_client.RestClient `inject:""`
}

func (this *RedirectApiRestController) InitController() {}

func (this *RedirectApiRestController) Redirect(ctx *gin.Context) {
	this.RedirectToService(ctx, this.RestClient)
}
