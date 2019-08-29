package utils

import "github.com/gin-gonic/gin"

// Mapped diagnostics context
type MDC struct {
	context map[string]string
}

func NewMDC() *MDC {
	return &MDC{
		context: map[string]string{},
	}
}

func (this *MDC) Ψ(key, value string) {
	this.context[key] = value
}

func (this *MDC) Put(ctx *gin.Context) {
	for k, v := range this.context {
		ctx.Header("$MDC$."+k, v)
	}

}
