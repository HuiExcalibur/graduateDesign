package controller

import (
	"Shaw/goWeb/chatRoom/data"
	"Shaw/goWeb/chatRoom/room"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	ExpireTime = time.Now().Add(72 * time.Hour)
	jwtKey     = []byte("chat room key")
)

type Claims struct {
	UserName string
	PassWord string
	jwt.StandardClaims
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//写入数据库
	db := data.GetDB()

	err := db.Register(username, password)
	if err != nil {
		log.Println("register failed ", err)
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})

		// 终止处理
		c.Abort()
		return
	}

	jwtToken := SetToken(username, password)
	if jwtToken != "" {
		c.SetCookie("jwt-token", jwtToken, int(time.Hour*72), "/", "127.0.0.1", false, false)
		c.SetCookie("user", username, int(time.Hour*72), "/", "127.0.0.1", false, false)
		c.SetCookie("nickname", username, int(time.Hour*72), "/", "127.0.0.1", false, false)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}

	c.JSON(http.StatusBadGateway, gin.H{
		"status": "failure",
	})
	// c.Redirect(http.StatusPermanentRedirect, "http://127.0.0.1:8080/index")
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	db := data.GetDB()
	nickname, err := db.Login(username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":   "failure",
			"error":    err.Error(),
			"username": username,
			"password": password,
		})
		// c.Next()

		c.Abort()
		return
	}

	jwtToken := SetToken(username, password)
	if jwtToken != "" {
		c.SetCookie("jwt-token", jwtToken, int(time.Hour*72), "/", "127.0.0.1", false, false)
		c.SetCookie("user", username, int(time.Hour*72), "/", "127.0.0.1", false, false)
		c.SetCookie("nickname", nickname, int(time.Hour*72), "/", "127.0.0.1", false, false)
		// c.JSON(http.StatusBadGateway, gin.H{
		// 	"status": "success",
		// })

		// c.Redirect(http.StatusPermanentRedirect, "http://127.0.0.1:8080/index")

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}

	c.JSON(http.StatusBadGateway, gin.H{
		"status": "failure",
	})
}

func NewRoom(c *gin.Context) {
	roonmane := c.Query("roomname")
	username := c.Query("username")

	db := data.GetDB()
	err := db.NewRoom(roonmane, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
		})
		c.Abort()
		return
	}
	//new room in the hub
	hub := room.GetHub()

	new_room := room.NewRoom(roonmane)
	member := hub.Users[username]

	hub.RegisterRoom <- new_room
	new_room.Register <- member

	//new room in the database

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func GetRoom(c *gin.Context) {
	username, err := c.Cookie("user")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"error":  err.Error(),
		})
		return
	}

	db := data.GetDB()

	ret, err := db.GetRoom(username)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"error":  err.Error(),
		})
		return
	}
	log.Println("get room success ", ret)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"rooms":  ret,
	})
}

func SearchRoom(c *gin.Context) {
	keyword := c.Query("key")

	db := data.GetDB()
	ret, err := db.SearchRoom(keyword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  err.Error(),
			"status": "failure",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"rooms":  ret,
	})
}

func EnterRoom(c *gin.Context) {
	roomname := c.Query("roomname")
	username, _ := c.Cookie("user")

	db := data.GetDB()
	err := db.EnterRoom(roomname, username)
	if err != nil {

		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
		})
		return
	}

	hub := room.GetHub()

	enterRoom := hub.Rooms[roomname]
	enterUser := hub.Users[username]

	enterRoom.Register <- enterUser

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func QuitRoom(c *gin.Context) {
	roomname := c.Query("roomname")
	username := c.Query("username")

	db := data.GetDB()
	err := db.QuitRoom(roomname, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
		})
		c.Abort()
		return
	}

	hub := room.GetHub()
	quitRoom := hub.Rooms[roomname]
	quitUser := hub.Users[username]
	// hub.UnRegisterRoom <- quitRoom
	quitRoom.UnRegister <- quitUser

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func History(c *gin.Context) {
	roomname := c.Query("roomname")
	username, _ := c.Cookie("user")

	db := data.GetDB()
	rows, err := db.History(roomname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error":  err.Error(),
			"status": "failure",
		})
		c.Abort()
		return
	}

	var msg room.Msg
	hub := room.GetHub()
	user := hub.Users[username]

	for rows.Next() {
		err := rows.Scan(&msg.Data, &msg.Username)
		if err != nil {
			return
		}
		msg.Room = roomname
		user.Send <- msg
		// hub.History <- msg
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func ChangeNickname(c *gin.Context) {
	userName, _ := c.Cookie("user")
	nickname := c.Query("nickname")

	db := data.GetDB()
	err := db.ChangeNickname(userName, nickname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "failure",
			"error":  err.Error(),
		})

		return
	}
	c.SetCookie("nickname", nickname, int(time.Hour*72), "/", "127.0.0.1", false, false)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func SetToken(username, password string) string {
	claim := &Claims{
		UserName: username,
		PassWord: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LuPengyi",
			Subject:   "登录鉴权",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return tokenStr
	// c.Header("Authorization", tokenStr)

	// c.SetCookie("jwt-token", tokenStr, 3600, "/", "127.0.0.1", false, false)
	// c.SetCookie("name", "xiaoming", 3600, "/", "127.0.0.1", false, false)

	// fmt.Println(tokenStr)
	// c.JSON(200, gin.H{
	// 	"token": tokenStr,
	// })
}

func ParseToken(tokenStr string) (*jwt.Token, *Claims, error) {
	claim := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claim, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return token, claim, err
}
