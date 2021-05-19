package main

import (
	"fmt"
	"net/http"

	"Shaw/goWeb/chatRoom/controller"
	"Shaw/goWeb/chatRoom/room"

	"github.com/gin-gonic/gin"
)

// var wsUpgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

var mainHub *room.Hub

// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsUpgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("failed to upgrade websocket", err)
// 		return
// 	}

// 	user := &room.User{
// 		Name:  conn.RemoteAddr().String(),
// 		Hub:   mainHub,
// 		Rooms: []string{"mainRoom"},
// 		Conn:  conn,
// 		Send:  make(chan room.Msg, 64),
// 	}
// 	fmt.Println(user.Name)
// 	user.Hub.RegisterUser <- user
// 	fmt.Println("has register", user.Name)
// 	go user.Read()
// 	go user.Write()

// 	// for {
// 	// 	msgtype, msg, err := conn.ReadMessage()
// 	// 	if err != nil {
// 	// 		fmt.Println("read failed ", err)
// 	// 		break
// 	// 	}
// 	// 	rpy := "received: " + string(msg) + " in " + time.Now().Format("15:04:05")
// 	// 	fmt.Println(rpy+" from ", conn.RemoteAddr())
// 	// 	conn.WriteMessage(msgtype, []byte(rpy))
// 	// }
// }

func test(c *gin.Context) {
	//c.HTML(200,)
	fmt.Println("get access from", c.Request.RemoteAddr)
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	r := gin.Default()

	mainHub = room.GetHub()
	go mainHub.Run()

	cRoom := room.NewRoom("mainRoom")
	mainHub.RegisterRoom <- cRoom

	r.LoadHTMLFiles("static/html/index.html", "static/html/login.html")
	// r.Static("/js", "./static/js")
	r.StaticFS("/public", http.Dir("./static"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)

	r.GET("/test", func(c *gin.Context) {
		username, _ := c.Cookie("user")
		token, _ := c.Cookie("jwt-token")
		c.JSON(200, gin.H{
			"user":      username,
			"jwt-token": token,
		})
	})
	r.GET("/getroom", controller.GetRoom)

	r.GET("/index", test)

	r.GET("/WS", func(c *gin.Context) {
		username, _ := c.Cookie("user")
		controller.WSHandler(c.Writer, c.Request, username)
	})

	r.Run(":8080")
}
