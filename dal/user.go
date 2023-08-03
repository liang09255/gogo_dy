package dal

import (
	"errors"
	"gorm.io/gorm"
	"main/global"
)

const MaxUsernameLength = 32

// TODO 字段不全，根据需要补充

type User struct {
	gorm.Model
	ID            int64  `gorm:"primarykey" json:"user_id"`
	Username      string `gorm:"not null" json:"user_name"`
	Password      string `gorm:"not null" json:"password"`
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
	IsFollow      bool   `gorm:"is_follow" json:"is_follow"`
}

type UserInfoResponse struct {
	ID            int64  `gorm:"primarykey" json:"id"`
	Username      string `gorm:"not null" json:"name"`
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
	IsFollow      bool   `gorm:"is_follow" json:"is_follow"`
}

type userDal struct{}

var UserDal = &userDal{}

func (u *userDal) CreateUser(user *User) error {
	//数据校验
	if user.Username == "" {
		return errors.New("用户名为空")
	}
	if len(user.Username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if user.Password == "" {
		return errors.New("密码为空")
	}
	t := global.MysqlDB.Create(user).First(&user)
	return t.Error
}

func (u *userDal) GetUserInfoById(userId int64, infoResponse *UserInfoResponse) error {
	var user = new(User)
	//直接传infoResponse会查询错误的数据表
	global.MysqlDB.Where("id=?", userId).First(&user)
	//.Select([]string{"id", "username", "follow_count", "follower_count", "is_follow"})
	infoResponse.FollowCount = user.FollowCount
	infoResponse.FollowerCount = user.FollowerCount
	infoResponse.ID = user.ID
	infoResponse.IsFollow = user.IsFollow
	infoResponse.Username = user.Username
	//id为零值，说明sql执行失败
	if infoResponse.ID == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}
