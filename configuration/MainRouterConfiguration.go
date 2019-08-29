package configuration

import limit "github.com/bu/gin-access-limit"
import "github.com/gin-gonic/gin"
import "github.com/golang/glog"
import "github.com/szuecs/gin-glog"
import "github.com/szuecs/gin-gomonitor"
import "github.com/szuecs/gin-gomonitor/aspects"
import "gopkg.in/mcuadros/go-monitor.v1/aspects"
import "time"

var Router = gin.Default()

func MainRouterConfiguration() {
	// programatically set swagger info
	//docs.SwaggerInfo.Title = "Swagger Example API"
	//docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	//docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = "petstore.swagger.io"
	//docs.SwaggerInfo.BasePath = "/v2"

	//Router := gin.Default()

	//Router.Use(limit.CIDR("172.18.0.0/16"))
	//security_Configuration()

	logger := ginglog.Logger(3 * time.Second)
	Router.Use(func(ctx *gin.Context) {
		logger(ctx)
	})
	Router.Use(gin.Recovery())

	glog.Warning("warning")
	glog.Error("err")
	glog.Info("info")
	glog.V(2).Infoln("This line will be printed if you use -v=N with N >= 2.")

	//AspectSetup()

}

func AspectSetup() {
	counterAspect := ginmon.NewCounterAspect()
	counterAspect.StartTimer(1 * time.Minute)
	anotherAspect := &CustomAspect{3}
	asps := []aspects.Aspect{counterAspect, anotherAspect}
	Router.Use(ginmon.CounterHandler(counterAspect))
	gomonitor.Start(9000, asps)
}

func security_Configuration() gin.IRoutes {
	return Router.Use(limit.CIDR("127.0.0.1/16"))
}

type CustomAspect struct {
	CustomValue int
}

func (a *CustomAspect) GetStats() interface{} {
	return a.CustomValue
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}
