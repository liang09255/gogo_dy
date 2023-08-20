package dal

import (
	"api/global"
	"errors"
	"gorm.io/gorm"
)

const MaxUsernameLength = 32

type User struct {
	gorm.Model
	ID             int64  `gorm:"primarykey" json:"user_id"`
	Username       string `gorm:"not null" json:"user_name"`
	Password       string `gorm:"not null" json:"password"`
	FollowCount    int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount  int64  `gorm:"default:0" json:"follower_count"`
	TotalFavorited int64  `gorm:"default:0" json:"total_favorited"` // 获点赞数量
	WorkCount      int64  `gorm:"default:0" json:"work_count"`      // 作品数
	FavoriteCount  int64  `gorm:"default:0" json:"favorite_count"`  //喜欢数
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

func (u *userDal) GetUserInfoById(userId int64) (user *User, err error) {
	t := global.MysqlDB.Where("id=?", userId).First(&user)
	if errors.Is(t.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}

	if t.Error != nil {
		return nil, t.Error
	}

	return user, nil
}

func (u *userDal) GetUserByUserName(userName string, user *User) error {
	//直接传infoResponse会查询错误的数据表
	t := global.MysqlDB.Where("username = ?", userName).Find(user)
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

func (u *userDal) MGetUser(ids []int64) (users []*User, err error) {
	t := global.MysqlDB.Where("id in ?", ids).Find(&users)
	if t.Error != nil {
		return nil, t.Error
	}
	return users, nil
}

func (u *userDal) AddFollowCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
}

func (u *userDal) SubFollowCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
}

func (u *userDal) AddFollowerCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
}

func (u *userDal) SubFollowerCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
}

func (u *userDal) AddWorkCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("work_count", gorm.Expr("work_count+?", 1)).Error
}

func (u *userDal) SubWorkCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("work_count", gorm.Expr("work_count - ?", 1)).Error
}

func (u *userDal) AddTotalFavorited(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("total_favorited", gorm.Expr("total_favorited + ?", 1)).Error
}

func (u *userDal) SubTotalFavorited(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("total_favorited", gorm.Expr("total_favorited - ?", 1)).Error
}

func (u *userDal) AddFavoriteCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (u *userDal) SubFavoriteCount(uid int64) error {
	return global.MysqlDB.Model(&User{}).Where("id = ?", uid).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}
