package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lthnh15032001/ngrok-impl/docs"
	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/api/controllers"
	"github.com/lthnh15032001/ngrok-impl/internal/store"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Init(errChan chan error) (bool error) {
	var err error
	var scInterface store.Interface
	config := config.GetConfig()
	scInterface, err = store.GetOnce()
	if err != nil {
		errChan <- err
	}
	router := gin.New()
	// gin.SetMode("")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//Routes for healthcheck of api server
	healthcheck := router.Group("health")
	{
		health := new(controllers.HealthController)
		ping := new(controllers.PingController)
		healthcheck.GET("/health", health.Status)
		healthcheck.GET("/ping", ping.Ping)
	}

	//Routes for swagger
	swagger := router.Group("swagger")
	{
		// programatically set swagger info
		docs.SwaggerInfo.Title = "Iot Rogo Stream REST API "
		docs.SwaggerInfo.Description = "This is a backend Rogo API."
		docs.SwaggerInfo.Version = "1.0"
		// docs.SwaggerInfo.Host = "cloudfactory.swagger.io"
		// docs.SwaggerInfo.BasePath = "/v1"

		swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	user := router.Group("user")
	{
		userController := &controllers.UserController{
			StoreInterface: scInterface,
		}

		user.GET("/", userController.GetAllUsers)
		user.POST("/", userController.GetAllUsers)
	}

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	errRunCfg := router.Run(config.GetString("server.port"))
	if errRunCfg != nil {
		errChan <- err
	}
	return err
}
