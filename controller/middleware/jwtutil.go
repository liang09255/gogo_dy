package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"main/dal"
	"main/global"
	"time"
)

var jwtKey []byte

func jwtUtilInit() {
	jwtKey = []byte(global.Config.JwtKey)
}

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

func ReleaseToken(user dal.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gogo_dy_auth",
			Subject:   "user_token",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "jwt generation error", err
	}
	return tokenString, nil
}
