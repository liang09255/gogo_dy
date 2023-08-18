package user_service_v1

import (
	"common/idl/user"
	"context"
	"fmt"

	"time"
	"user/internal/dao"
	"user/internal/database"
	"user/internal/database/tran"
	"user/internal/domain"
	"user/internal/repo"
)

type UserService struct {
	user.UnimplementedUserServiceServer
	cache       repo.Cache
	userDomain  *domain.UserDomain
	transaction tran.Transaction
}

func New() *UserService {
	return &UserService{
		cache:      dao.Rc,
		userDomain: domain.NewUserDomain(),
	}
}

func (us *UserService) Register(ctx context.Context, msg *user.RegisterRequest) (*user.RegisterResponse, error) {

	// 校验用户是否已存在
	// 输入：用户名
	// 返回：是否存在的布尔值
	userName := msg.Username
	passWord := msg.Password

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	NotExist, err := us.userDomain.UserIsNotExist(c, userName)
	if err != nil {
		return &user.RegisterResponse{}, err
	}

	if NotExist == false {
		err = fmt.Errorf("用户已存在")
		return &user.RegisterResponse{}, err
	}

	// 密码加密
	// passWord = encrypts.Md5(passWord + global.Config.PasswordSalt)

	// 插入数据库
	err = us.transaction.Action(func(conn database.DbConn) error {
		err = us.userDomain.AddUser(conn, c, userName, passWord)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &user.RegisterResponse{}, err
	}

	////jwt分发和返回值封装
	//var user = new(dal.User)
	//err = dal.UserDal.GetUserByUserName(userName, user)
	//if err != nil {
	//	return
	//}
	//
	//token, err := middleware.ReleaseToken(user.ID)
	//if err != nil {
	//	err = fmt.Errorf("token分发失败")
	//	return
	//}

	// 封装返回结果
	return &user.RegisterResponse{UserId: 1, Token: "token"}, nil
}
