package model

type UserRegisterRequest struct {
	Password string // 密码，最长32个字符
	Username string // 注册用户名，最长32个字符
}

type UserRegisterResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

type UserLoginRequest struct {
	Password string // 密码，最长32个字符
	Username string // 注册用户名，最长32个字符
}

type UserLoginResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	Token      *string `json:"token"`       // 用户鉴权token
	UserID     *int64  `json:"user_id"`     // 用户id
}

type UserInfoRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type UserInfoResponse struct {
	StatusCode int64     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string   `json:"status_msg"`  // 返回状态描述
	User       *Userinfo `json:"user"`        // 用户信息
}
