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
		c.HTML(http.StatusOK, "index.html", "")
	})
	router.POST("login", handler.Authenticated)
	router.POST("logout", handler.AuthLogout)

	router.GET("/history", handler.ReadHistory)
	send := router.Group("/send", gin.BasicAuth(gin.Accounts{"dev": "2104"}))
	{
		send.POST("wa", handler.SendMessage)
	}

	return router
}
