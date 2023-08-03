package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"main/dal"
	"time"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

func ReleaseToken(user dal.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "gogo_dy",
			Subject:   "user token",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "jwt generation error", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	//claim里可以拿到userId
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}

//TODO jwt鉴权中间件
