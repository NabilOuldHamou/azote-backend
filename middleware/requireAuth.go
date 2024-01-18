package middleware

import (
	"azote-backend/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RequireAuth(c *gin.Context) {

	session, err := token.ParseToken(c)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if time.Now().Unix() > session.ExpiresAt.Unix() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userId", session.Bearer.String())
	c.Next()
}
