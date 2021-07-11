package main

import (
	"fmt"
	"net/http"

	"Shaw/goWeb/chatRoom/controller"
	"Shaw/goWeb/chatRoom/midware"
	"Shaw/goWeb/chatRoom/room"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	//c.HTML(200,)
	fmt.Println("get access from", c.Request.RemoteAddr)
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	r := gin.Default()

	room.GetHub()

	r.LoadHTMLFiles("static/html/index.html", "static/html/login.html")
	// r.Static("/js", "./static/js")
	r.StaticFS("/public", http.Dir("./static"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)

	authorize := r.Group("/")

	authorize.Use(midware.Auth())
	{
		authorize.GET("/changenickname", controller.ChangeNickname)
		authorize.GET("/newroom", controller.NewRoom)
		authorize.GET("/search", controller.SearchRoom)
		authorize.GET("/quit", controller.QuitRoom)
		authorize.GET("/enter", controller.EnterRoom)
		authorize.GET("/getroom", controller.GetRoom)
		authorize.GET("/history", controller.History)
		authorize.GET("/index", index)
		authorize.GET("/WS", func(c *gin.Context) {
			username, _ := c.Cookie("user")
			controller.WSHandler(c.Writer, c.Request, username)
		})
	}

	r.Run(":8080")
}
