package dal

import (
	"fmt"
	"main/global"
	"main/middleware"
	"main/model"
	"strconv"
)

type userDal struct{}

var UserDal = &userDal{}

// Exist 判断用户是否存在
func (h *userDal) Exist(username string) error {
	var db = global.MysqlDB
	var user = model.User{Username: username}
	t := db.First(&user)
	return t.Error
}

// Adduser 增加新用户
func (h *userDal) Adduser(username string, password string) error {
	var db = global.MysqlDB
	var user = model.User{
		Username:      username,
		Password:      password,
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
	}
	t := db.Create(&user)
	return t.Error
}

// GetUserID 获取用户ID
func (h *userDal) GetUserID(username string, password string) (int64, error) {
	var db = global.MysqlDB
	var res model.User

	t := db.Where("username = ? AND password = ?", username, password).First(&res)

	return res.Id, t.Error
}

// GetUserIDbyname 通过username获取用户ID
func (h *userDal) GetUserIDbyname(username string) (int64, error) {
	var db = global.MysqlDB
	var res model.User

	t := db.Where("username = ? ", username).First(&res)

	return res.Id, t.Error
}

// GetUserInfo 获取用户信息 参数包括：当前用户token 待获取用户信息
func (h *userDal) GetUserInfo(token string, userid string, resp *model.Userinfo) error {
	var db = global.MysqlDB
	var res model.User

	fmt.Println(userid)

	id, err := strconv.Atoi(userid)

	if err != nil {
		return err
	}

	t := db.Where("id = ?", id).First(&res)

	resp.ID = res.Id
	resp.Name = res.Name
	resp.FollowerCount = res.FollowerCount
	resp.FollowCount = res.FollowCount
	resp.FavoriteCount = 0
	resp.WorkCount = 0
	resp.Signature = "signature"
	resp.BackgroundImage = "background"
	resp.Avatar = "avatar"
	resp.TotalFavorited = "0"

	//需要查看当前token来判断当前用户
	CurrentUserName := middleware.GetUsernameByToken(token)

	CurrentId, _ := UserDal.GetUserIDbyname(CurrentUserName)

	fg, _ := RelationDal.CheckRelation(CurrentId, int64(id))

	if fg == true {
		resp.IsFollow = true
	} else {
		resp.IsFollow = false
	}

	return t.Error
}
