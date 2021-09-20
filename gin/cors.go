package d_gin

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetCors(router *gin.Engine, allowHeaders []string)  {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, allowHeaders...)
	router.Use(cors.New(config))
	return
}
