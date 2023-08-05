package model

type DouyinRelationActionRequest struct {
	ActionType string `json:"action_type"` // 1-关注，2-取消关注
	ToUserID   string `json:"to_user_id"`  // 对方用户id
	Token      string `json:"token"`       // 用户鉴权token
}

type DouyinRelationActionResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type DouyinRelationFollowerListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFollowerListResponse struct {
	StatusCode string     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string    `json:"status_msg"`  // 返回状态描述
	UserList   []Userinfo `json:"user_list"`   // 用户列表
}

type DouyinRelationFollowListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFollowListResponse struct {
	StatusCode string     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string    `json:"status_msg"`  // 返回状态描述
	UserList   []Userinfo `json:"user_list"`   // 用户信息列表
}

type DouyinRelationFriendListRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type DouyinRelationFriendListResponse struct {
	StatusCode string     `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string    `json:"status_msg"`  // 返回状态描述
	UserList   []Userinfo `json:"user_list"`   // 用户列表
}
