package token

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"math"
	"os"
	"time"
)

type UserSession struct {
	Bearer    uuid.UUID
	ExpiresAt time.Time
}

func ParseToken(c *gin.Context) (UserSession, error) {
	tokenString, err := c.Cookie("Authorization")
	if len(tokenString) == 0 || err != nil {
		return UserSession{}, errors.New("cookie not found")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user models.User
		initializers.DB.First(&user, "id = ?", claims["bearer"])

		sec, dec := math.Modf(claims["expiresAt"].(float64))

		return UserSession{Bearer: user.ID, ExpiresAt: time.Unix(int64(sec), int64(dec*(1e9)))}, nil

	} else {
		return UserSession{}, errors.New("token is not valid")
	}
}
