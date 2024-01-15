package middleware

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"azote-backend/tokens"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	var user models.User
	initializers.DB.First(&user, "id = ?", session.Bearer)
	if user.ID == uuid.Nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", user)
	c.Next()
}
