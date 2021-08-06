package routers

import (
	"github.com/aldinokemal/go-whatsapp-web/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routers() *gin.Engine {
	router := gin.Default()
	router.Static("/statics", "./statics")    // static path
	router.LoadHTMLGlob("./templates/*.html") // template

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/app/home")
	})

	//app := router.Group("app", gin.BasicAuth(gin.Accounts{config.BasicAuthUser: config.BasicAuthPswd}))
	app := router.Group("app")
	{
		app.GET("home", handler.Home)
		app.GET("home/get-account", handler.GetAccount)

		app.POST("login", handler.Authenticated)
		app.POST("logout", handler.AuthLogout)

		app.POST("send/text", handler.SendMessage)
		app.POST("send/image", handler.SendImage)

		router.GET("/history", handler.ReadHistory)
	}

	return router
}
