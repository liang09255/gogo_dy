package middleware

import (
	"strings"
)

func GetUsernameByToken(token string) string {
	//根据token获取用户名
	// token = username & password
	s := strings.Split(token, "&")
	name := s[0]
	return name
}
