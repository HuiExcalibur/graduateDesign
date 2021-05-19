package midware

import (
	"Shaw/goWeb/chatRoom/controller"
	"log"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("jwt-token")
		if err != nil {
			log.Println("get cookie fail")
			c.Abort()
		}

		token, _, err := controller.ParseToken(tokenStr)
		if !token.Valid || err != nil {
			log.Println("token is invalid")
			c.Abort()
		}
		c.Next()
	}
}
