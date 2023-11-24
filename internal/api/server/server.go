package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lthnh15032001/ngrok-impl/docs"
	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/api/controllers"
	middlewares "github.com/lthnh15032001/ngrok-impl/internal/api/middlewares"
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
	router := gin.Default()

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}

	// router.Use(cors.New(configCors))
	router.Use(CORSMiddleware())

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

		user.GET("/", middlewares.AuthMiddleware("getUsers"), userController.GetAllUsers)
		user.POST("/", userController.GetAllUsers)
	}

	tunnel := router.Group("tunnel")
	{
		tunnelController := &controllers.TunnelController{
			StoreInterface: scInterface,
		}

		tunnel.GET("/", middlewares.AuthMiddleware("get-tunnel"), tunnelController.GetTunnelActive)
		tunnel.POST("/", middlewares.AuthMiddleware("get-tunnel"), tunnelController.AddTunnel)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Not Found API",
		})
		// c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	errRunCfg := router.Run(config.GetString("server.port"))
	if errRunCfg != nil {
		errChan <- err
	}
	return err
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
