package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaims struct {
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	jwt.StandardClaims
}

func CreateJWT(phone_number string, status string, expired int64) string {
	claims := &jwtCustomClaims{
		phone_number,
		status,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err.Error()
	}
	return t
}
