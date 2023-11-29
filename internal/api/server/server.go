package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/lthnh15032001/ngrok-impl/docs"
	"github.com/lthnh15032001/ngrok-impl/internal/api/config"
	"github.com/lthnh15032001/ngrok-impl/internal/api/controllers"
	middlewares "github.com/lthnh15032001/ngrok-impl/internal/api/middlewares"
	"github.com/lthnh15032001/ngrok-impl/internal/store"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"
)

func Init(errChan chan error) (bool error) {
	var err error
	var db *gorm.DB
	var scInterface store.Interface
	config := config.GetConfig()
	scInterface, db, err = store.GetOnce()
	gs := gormsessions.NewStore(db, true, []byte("secret"))
	if err != nil {
		errChan <- err
	}

	router := gin.Default()
	router.Use(sessions.Sessions("api-session-storage", gs))
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

		user.GET("/", middlewares.AuthMiddleware("get-user"), userController.GetAllUsers)
		user.GET("/:id", middlewares.AuthMiddleware("get-user"), userController.GetUser)
		user.POST("/", middlewares.AuthMiddleware("add-user"), userController.AddUserAuthen)
		user.PATCH("/:id", middlewares.AuthMiddleware("edit-user"), userController.EditUserAuthen)
		user.DELETE("/", middlewares.AuthMiddleware("delete-user"), userController.DeleteUserAuthen)
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
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
