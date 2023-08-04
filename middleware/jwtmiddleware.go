package middleware

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"log"
	"main/dal"
	"net/http"
	"time"
)

// biz/router/middleware/jwtutil.go

// Claim 定义用户登陆信息结构体
type Claim struct {
	ID       int64
	Username string
}

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "user_id"
)

func JwtMwInit() {
	var userId int64
	// the jwt middleware
	JwtMiddleware1, err := jwt.New(&jwt.HertzJWTMiddleware{
		// 置所属领域名称
		Realm: "hertz jwt",
		// 用于设置签名密钥
		Key: []byte("secret_key&&gogo_dy"),
		// 设置 token 过期时间
		Timeout: time.Hour * 24,
		// 设置最大 token 刷新时间
		MaxRefresh: time.Hour * 4,
		// 设置 token 的获取源
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// 设置从 header 中获取 token 时的前缀
		TokenHeadName: "Bearer",
		// 用于设置检索身份的键
		IdentityKey: IdentityKey,

		// 从 token 提取用户信息
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return claims[IdentityKey]
		},
		// 认证
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 32); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			user, err := dal.UserDal.CheckUser(loginStruct.Username, loginStruct.Password)
			if err != nil {
				return nil, err
			}
			if user == nil {
				return nil, errors.New("user already exists or wrong password")
			}
			c.Set("user_id", user.ID)
			userId = user.ID
			// 设置jwt负载的信息
			return user.ID, err
		},
		// 登陆成功后为向 token 中添加自定义负载信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		// 登录校验成功，将token返回给前端
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"status_code ": code,
				"status_msg ":  "登陆成功",
				"user_id ":     userId,
				"token":        token,
			})
			//controller.LoginSuccessResponse(c, "success", LoginResponse)
			//hlog.CtxInfof(ctx, "Login success ，token is issued clientIP: "+c.ClientIP())
			//c.Set("token", token)
		},
		// 鉴权
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			// 单一角色 不设权限校验
			if v, ok := data.(float64); ok {
				current_user_id := int64(v)
				c.Set("current_user_id", current_user_id)
				hlog.CtxInfof(ctx, "Token is verified clientIP: "+c.ClientIP())
				return true
			}
			return false
		},
		// 设置 jwt 校验流程发生错误时响应所包含的错误信息
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		// jwt 验证流程失败的响应函数
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			BaseFailResponse(c, message)
		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	JwtMiddleware = JwtMiddleware1
}

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func BaseFailResponse(ctx *app.RequestContext, msg string) {
	Response(ctx, BaseResponse{StatusCode: 1, StatusMsg: msg})
}
func Response(ctx *app.RequestContext, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}
