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
}

type UserInfoResponse struct {
	ID            int64  `gorm:"primarykey" json:"id"`
	Username      string `gorm:"not null" json:"name"`
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
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

func (u *userDal) IsNotExist(userName string) bool {
	user := new(User)
	global.MysqlDB.Where("username = ?", userName).First(&user)
	if user.ID != 0 {
		return false
	}
	return true
}

func (u *userDal) GetUserInfoById(userId int64, infoResponse *UserInfoResponse) error {
	var user = new(User)
	//直接传infoResponse会查询错误的数据表
	t := global.MysqlDB.Where("id=?", userId).First(&user)
	//.Select([]string{"id", "username", "follow_count", "follower_count", "is_follow"})
	infoResponse.FollowCount = user.FollowCount
	infoResponse.FollowerCount = user.FollowerCount
	infoResponse.ID = user.ID
	infoResponse.Username = user.Username
	//id为零值，说明sql执行失败
	if infoResponse.ID == 0 {
		return errors.New("该用户不存在")
	}
	return t.Error
}

func (u *userDal) GetUserByUserName(userName string, user *User) error {
	//直接传infoResponse会查询错误的数据表
	t := global.MysqlDB.Where("username = ?", userName).Find(&user)
	//id为零值，说明sql执行失败
	if user.ID == 0 {
		return errors.New("该用户不存在")
	}
	return t.Error
}

func (u *userDal) CheckUser(userName string, passWord string) (*User, error) {
	var user = new(User)
	t := global.MysqlDB.Where("username = ?", userName).Find(&user)
	//id为零值，说明sql执行失败
	if user.ID == 0 {
		return nil, errors.New("该用户不存在")
	}
	if user.Password != passWord {
		return nil, errors.New("用户验证失败")
	}
	return user, t.Error
}
