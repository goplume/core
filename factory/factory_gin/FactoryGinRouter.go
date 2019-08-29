package factory_gin

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/logger"
	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
	"log"
	"time"
)

type FactoryGinRouter struct {
	Log *logger.Logger
}

func (this *FactoryGinRouter) InitFactory() {
}

func (this *FactoryGinRouter) CreateGinRouter(
	serviceName string,
) (router *gin.Engine, routerGroup *gin.RouterGroup, address string) {

	config := configuration.NewServiceConfiguration(serviceName, "router.%s", this.Log)
	this.Log.RLog.Info("Read configuration from context " + config.ConfigurationContext)

	contextPath := config.GetString("context-path")
	address = config.GetString("address")

	router = gin.Default()
	router.Use(ginglog.Logger(3 * time.Second))
	//router.Use(this.Logger())
	router.Use(gin.Recovery())
	// Install nice.Recovery, passing the handler to call after recovery
	router.Use(nice.Recovery(recoveryHandler))

	routerGroup = router.Group(contextPath)
	gin.DisableConsoleColor()
	return
}

func recoveryHandler(c *gin.Context, err interface{}) {
	c.HTML(500, "error.tmpl", gin.H{
		"title": "Error",
		"err":   err,
	})
}

func (this *FactoryGinRouter) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawBody, _ := c.GetRawData()
		msg :=
			"GIN REQUEST BODY: =============================\n" +
				string(rawBody) +
				"-------------------------------------------------\n"
		this.Log.RLog.Info(msg)
		log.Println(msg)
	}
}
