package domain

import (
	"common/ggCompare"
	"common/ggEncrypts"
	"common/ggIDL/user"
	"common/ggLog"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	"time"
	"user/internal/dal"
	"user/internal/model"
	"user/internal/repo"
)

type UserDomain struct {
	userRepo      repo.UserRepo
	tranRepo      repo.TranRepo
	userCacheRepo repo.UserCacheRepo
}

func NewUserDomain() *UserDomain {
	return &UserDomain{
		userRepo:      dal.NewUserDao(),
		tranRepo:      dal.NewTranRepo(),
		userCacheRepo: dal.NewUserCacheRepo(),
	}
}

func (ud *UserDomain) Register(ctx context.Context, userName, passWord string) (userid int64, err error) {
	// 检查用户是否存在
	exist := ud.userRepo.Exist(ctx, userName)
	if exist {
		return 0, fmt.Errorf("用户已存在")
	}
	// 密码加密
	passWord = ggEncrypts.MD5Password(passWord)
	// 写入数据库
	var user = &model.User{
		Username: userName,
		Password: passWord,
	}
	err = ud.userRepo.CreateUser(ctx, user)
	return user.ID, err
}

func (ud *UserDomain) Login(ctx context.Context, userName, password string) (userid int64, err error) {
	password = ggEncrypts.MD5Password(password)
	var userDaoModel = &model.User{
		Username: userName,
		Password: password,
	}
	err = ud.userRepo.CheckUser(ctx, userDaoModel)
	return userDaoModel.ID, err
}

func (ud *UserDomain) MGetUserInfo(ctx context.Context, uids []int64) (userInfo []*user.UserInfoModel, err error) {
	var users []model.User
	// 先尝试从缓存中获取
	users, err = ud.userCacheRepo.MGetUserInfo(ctx, uids)
	// 检查是否有未命中缓存的
	var cacheUids []int64
	var missUids []int64
	for _, u := range users {
		cacheUids = append(cacheUids, u.ID)
	}
	ggLog.Debugf("cacheUids:%v", cacheUids)
	missUids = ggCompare.Diff(uids, cacheUids)
	// 从数据库中获取
	if len(missUids) != 0 {
		missUsers, err := ud.userRepo.MGetUserInfo(ctx, uids)
		if err != nil {
			return nil, err
		}
		users = append(users, missUsers...)
		// 异步写缓存
		gopool.Go(func() {
			if err := ud.userCacheRepo.MSetUserInfo(context.Background(), missUsers, 1*time.Minute); err != nil {
				ggLog.Error("写入缓存失败:", err)
			}
			ggLog.Debugf("写入缓存成功, uid:%v", missUids)
		})
	}
	// 封装返回值
	for _, u := range users {
		userInfo = append(userInfo, &user.UserInfoModel{
			Id:              u.ID,
			Name:            u.Username,
			FollowCount:     u.FollowCount,
			FollowerCount:   u.FollowerCount,
			Avatar:          "",
			BackgroundImage: "",
			Signature:       "",
			TotalFavorited:  u.TotalFavorited,
			WorkCount:       u.WorkCount,
			FavoriteCount:   u.FavoriteCount,
		})
	}
	return userInfo, nil
}
