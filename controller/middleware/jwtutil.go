package middleware

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dgrijalva/jwt-go"
	"main/global"
	"time"
)

var jwtKey []byte

func jwtUtilInit() {
	jwtKey = []byte(global.Config.JwtKey)
}

type claims struct {
	UserId int64 `json:"user_id"`
	jwt.StandardClaims
}

func ReleaseToken(userId int64) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &claims{
		UserId: userId,
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

func GetUserID(ctx *app.RequestContext) int64 {
	id, ok := ctx.Get(UserIDKey)
	if !ok {
		hlog.Error("get userID from context failed, path: %s", ctx.FullPath())
		ctx.Abort()
		return 0
	}
	return int64(id.(float64))
}
