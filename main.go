package main

import (
	"fmt"
	"net/http"

	"Shaw/goWeb/chatRoom/room"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mainHub *room.Hub

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("failed to upgrade websocket", err)
		return
	}

	user := &room.User{
		Addr:  conn.RemoteAddr().String(),
		Hub:   mainHub,
		Rooms: []string{"mainRoom"},
		Conn:  conn,
		Send:  make(chan room.Msg, 64),
	}
	fmt.Println(user.Addr)
	user.Hub.RegisterUser <- user
	fmt.Println("has register", user.Addr)
	go user.Read()
	go user.Write()

	// for {
	// 	msgtype, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("read failed ", err)
	// 		break
	// 	}
	// 	rpy := "received: " + string(msg) + " in " + time.Now().Format("15:04:05")
	// 	fmt.Println(rpy+" from ", conn.RemoteAddr())
	// 	conn.WriteMessage(msgtype, []byte(rpy))
	// }
}

func test(c *gin.Context) {
	//c.HTML(200,)
	fmt.Println("get access from", c.Request.RemoteAddr)
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	r := gin.Default()

	mainHub = room.NewHub()
	go mainHub.Run()

	cRoom := room.NewRoom("mainRoom")
	mainHub.RegisterRoom <- cRoom

	r.LoadHTMLFiles("static/html/index.html")
	// r.Static("/js", "./static/js")
	r.StaticFS("/public", http.Dir("./static"))

	r.GET("/", test)
	r.GET("/WS", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})

	r.Run(":8080")
}
