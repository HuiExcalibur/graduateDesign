package midware

import (
	"Shaw/goWeb/chatRoom/controller"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("jwt-token")
		if err != nil {
			log.Println("get cookie fail")
			// c.Request.URL.Path="/"

			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()
			return
		}

		token, _, err := controller.ParseToken(tokenStr)
		if !token.Valid || err != nil {
			log.Println("token is invalid")
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
