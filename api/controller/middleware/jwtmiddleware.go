package middleware

import (
	"context"
	"errors"
	"log"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/userCtlModel"
	"main/utils/encrypts"
	"time"

	"main/controller/ctlFunc"
	"main/dal"
	"main/global"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

const identityKey = "user_id"
const UserIDKey = "user_id"

var Jwt *jwt.HertzJWTMiddleware

func jwtMwInit() {
	var userId int64
	// the jwt middleware
	JwtMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		// 置所属领域名称
		Realm: "gogo_dy_auth",
		// 用于设置签名密钥
		Key: []byte(global.Config.JwtKey),
		// 设置 token 过期时间
		Timeout: time.Hour * 24,
		// 设置最大 token 刷新时间
		MaxRefresh: time.Hour * 4,
		// 设置 token 的获取源
		TokenLookup: "header: Authorization, query: token, cookie: jwt, form: token",
		// 设置从 header 中获取 token 时的前缀
		TokenHeadName: "Bearer",
		// 用于设置检索身份的键
		IdentityKey: identityKey,
		// 从 token 提取用户信息
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return claims[identityKey]
		},
		// 认证
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var reqObj userCtlModel.LoginReq
			if err := c.BindAndValidate(&reqObj); err != nil {
				return nil, err
			}

			username := reqObj.Username
			password := encrypts.Md5(reqObj.Password + global.Config.PasswordSalt)

			user, err := dal.UserDal.CheckUser(username, password)
			if err != nil {
				return nil, err
			}
			if user == nil {
				return nil, errors.New("user already exists or wrong password")
			}
			c.Set(UserIDKey, user.ID)
			userId = user.ID
			// 设置jwt负载的信息
			return user.ID, err
		},
		// 登陆成功后为向 token 中添加自定义负载信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		// 登录校验成功，将token返回给前端
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			var resp = &userCtlModel.LoginResp{
				BaseResp: baseCtlModel.NewBaseSuccessResp(),
				LoginResponse: userCtlModel.LoginResponse{
					UserId: userId,
					Token:  token,
				},
			}
			ctlFunc.Response(c, resp)
		},
		// 鉴权
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			// 单一角色 不设权限校验
			if v, ok := data.(float64); ok {
				currentUserId := int64(v)
				c.Set("current_user_id", currentUserId)
				return true
			}
			return false
		},
		// 设置 jwt 校验流程发生错误时响应所包含的错误信息
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			return e.Error()
		},
		// jwt 验证流程失败的响应函数
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			ctlFunc.BaseFailedResp(c, message)
		},
	})
	if err != nil {
		log.Fatal("JWT Init Error:" + err.Error())
	}

	Jwt = JwtMiddleware
}
