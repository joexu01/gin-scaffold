package router

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joexu01/gin-scaffold/docs"
	"github.com/joexu01/gin-scaffold/lib"
	"github.com/joexu01/gin-scaffold/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Example API
// @version         1.0
// @description     This is a sample server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8880
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//非登陆接口
	store := sessions.NewCookieStore([]byte("secret"))
	apiNormalGroup := router.Group("/api")
	apiNormalGroup.Use(sessions.Sessions("gin-session", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware())
	{
		// controller.ApiRegister(apiNormalGroup)
	}

	//登陆接口
	apiAuthGroup := router.Group("/api")
	apiAuthGroup.Use(
		sessions.Sessions("gin-session", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		//controller.ApiLoginRegister(apiAuthGroup)
	}

	return router
}
