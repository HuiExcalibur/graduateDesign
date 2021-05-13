package controller

import (
	"Shaw/goWeb/chatRoom/data"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	ExpireTime = time.Now().Add(60 * time.Second)
	jwtKey     = []byte("this is a key")
)

type Claims struct {
	UserID   string
	PassWord string
	jwt.StandardClaims
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	db, err := data.NewDB()
	if err != nil {
		return
	}

	db.Register(username, password)
}

func Login(c *gin.Context) {

}

func NewRoom(c *gin.Context) {

}

func DelRoom(c *gin.Context) {

}

func History(c *gin.Context) {

}
