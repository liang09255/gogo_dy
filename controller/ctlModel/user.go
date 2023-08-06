package ctlModel

type UserRegisterReq struct {
	Username string `query:"username,required" vd:"len($)>2 && len($)<32"`
	Password string `query:"password,required" vd:"len($)>5 && len($)<32"`
}

type UserLoginReq struct {
	Username string `query:"username,required" vd:"len($)>2 && len($)<32"`
	Password string `query:"password,required" vd:"len($)>5 && len($)<32"`
}

type UserLoginResp struct {
	BaseResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserInfoReq struct {
	UserId int64 `query:"user_id,required" vd:"$>0"`
}

func UserLoginSuccessResp() {

}
